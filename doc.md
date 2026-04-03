# devtunnel — Command Reference

Everything used to build, release, and maintain this project.

---

## Repository layout

```
devtunnel/
├── cmd/
│   ├── mytunnel/          # CLI binary (what users install as `mytunnel`)
│   ├── mytunneld/         # Tunnel server daemon (control + edge HTTP)
│   └── test-server/       # Sample Gin app for local testing
├── pkg/
│   └── tunnel/            # Public Go API for embedding (`import …/pkg/tunnel`)
├── internal/
│   ├── config/            # Shared defaults (ports, host suffix, dial address)
│   ├── id/                # Subdomain id generation
│   ├── protocol/          # Wire messages (TunnelRequest / TunnelResponse)
│   └── tunnel/
│       ├── client/        # Dial control plane, yamux, forward streams → local HTTP
│       └── server/        # Registry, TCP listener, edge HTTP handler
├── install.sh
├── go.mod
├── readme.md
└── doc.md
```

| Path | Role |
|------|------|
| `cmd/mytunnel` | CLI: connects from a developer machine and forwards a local port |
| `cmd/mytunneld` | Server: listens for tunnel clients and serves public HTTP by subdomain |
| `cmd/test-server` | Optional demo backend |
| `pkg/tunnel` | Stable import path for other Go modules (`tunnel.Start`, etc.) |
| `internal/tunnel/client` | Implementation: `Connect`, stream forwarding to `localhost` |
| `internal/tunnel/server` | Implementation: session registry, TCP acceptor, HTTP edge |
| `internal/protocol` | JSON line protocol between client streams and edge |
| `internal/config` | Single place for `:9000`, `:3000`, public host suffix |
| `internal/id` | Random subdomain tokens |

Only code under `pkg/` and `cmd/` entrypoints are intended as “public” surfaces. Everything under `internal/` is private to this module (Go enforces that).

---

## Naming

### Commands (`cmd/`)

| Name | Meaning |
|------|---------|
| **mytunnel** | The CLI users run on their machine to expose a local port (the shipped binary name). |
| **mytunneld** | The long-running **server** process. The trailing `d` follows common Unix daemon naming (`sshd`, `dockerd`). It listens on the control port (`:9000`) and the public HTTP edge (`:3000`). |

### Public library (`pkg/`)

| Name | Meaning |
|------|---------|
| **pkg/tunnel** | **`pkg/`** signals “safe for other repos to import.” The package name **tunnel** is a short, stable API for starting and stopping a tunnel, independent of marketing names. |

### Internal implementation (`internal/tunnel/`)

| Name | Meaning |
|------|---------|
| **internal/** | Only this module may import these packages; not a semver-stable API. |
| **tunnel** (directory) | Umbrella for all tunneling logic (both “client” and “server” sides of the product). |
| **client** (`internal/tunnel/client`) | The side that **dials** the control server and **forwards** yamux streams to local HTTP. “Client” means tunnel client, not “HTTP client” (though it uses `http.Client` internally). |
| **Connect** | Library-oriented entrypoint: open the tunnel and return the public URL plus a `stop` function. |
| **Options** | Optional settings (`ServerAddr`, `LocalHost`) without a long parameter list. |
| **forward.go** / **acceptStreams** / **handleStream** | Yamux multiplexes **streams**; each stream carries one request/response, **forwarded** to `http://localhost:<port>`. |
| **server** (`internal/tunnel/server`) | Accepts tunnel connections and routes public HTTP to the right session. |
| **Registry** | In-memory map: subdomain → yamux session (“who is connected right now”). |
| **TCPListener** | Listens on TCP for **control** connections (yamux). Named by transport so it is not confused with the HTTP edge. |
| **EdgeHTTP** | The **public HTTP** server that terminates browser traffic and routes by `Host` (subdomain). “Edge” = public entry point of the network path. |

### Shared packages

| Name | Meaning |
|------|---------|
| **internal/protocol** | Wire contract (JSON) between edge and tunnel client streams. |
| **internal/config** | Defaults so CLI, daemon, and library agree on ports and host patterns. |
| **internal/id** | Generates subdomain-safe random ids. |

### Mental model (one line each)

| Piece | In one sentence |
|-------|------------------|
| `mytunnel` | Run on my laptop to expose a local port through the tunnel server. |
| `mytunneld` | Run on the machine that owns the public hostname and listens for tunnel clients. |
| `pkg/tunnel` | Import in Go to start/stop a tunnel without using the CLI. |
| `internal/tunnel/client` | Dial the control plane and forward each stream to local HTTP. |
| `internal/tunnel/server` | Accept tunnels and map public HTTP requests to the correct session. |

If you align branding with the repo name (`devtunnel`), you could rename binaries to `devtunnel` / `devtunneld`; `pkg/tunnel` can stay or become `pkg/devtunnel` depending on how you want the import path to read.

---

## Go Build Commands

### Build without cache (`-a`)

```bash
GOOS=darwin GOARCH=arm64 go build -a -o mytunnel-mac-arm64 ./cmd/mytunnel
```

| Part | Meaning |
|------|---------|
| `GOOS=darwin` | Target OS: macOS |
| `GOARCH=arm64` | Target CPU: Apple Silicon (M1/M2/M3) |
| `go build` | Compile the Go program |
| `-a` | Force rebuild of all packages — skips the build cache |
| `-o mytunnel-mac-arm64` | Output binary filename |
| `./cmd/mytunnel` | CLI entry point |

### All three platform builds

```bash
# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -a -o mytunnel-mac-arm64 ./cmd/mytunnel

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -a -o mytunnel-mac ./cmd/mytunnel

# Linux x86_64
GOOS=linux GOARCH=amd64 go build -a -o mytunnel-linux ./cmd/mytunnel
```

### Full clean rebuild (nuclear option)

```bash
go clean -cache -modcache -i -r
go build -a -o mytunnel ./cmd/mytunnel
```

| Flag | Meaning |
|------|---------|
| `-cache` | Delete the build cache |
| `-modcache` | Delete the downloaded module cache |
| `-i` | Remove installed packages |
| `-r` | Apply recursively to all dependencies |

---

## Install Locally

```bash
sudo cp mytunnel-mac-arm64 /usr/local/bin/mytunnel
sudo chmod +x /usr/local/bin/mytunnel
```

| Command | Meaning |
|---------|---------|
| `sudo cp` | Copy with root privileges |
| `/usr/local/bin/` | Standard location for user-installed binaries — already in `$PATH` |
| `chmod +x` | Mark the file as executable |

---

## GitHub CLI (`gh`)

### Install

```bash
brew install gh
```

### Authenticate

```bash
gh auth login
```

Follow the prompts: GitHub.com → HTTPS → Login with a web browser.

### Upload binaries to a release

```bash
gh release upload devtunnel mytunnel-mac mytunnel-mac-arm64 mytunnel-linux --clobber
```

| Part | Meaning |
|------|---------|
| `gh release upload` | Upload assets to an existing GitHub release |
| `devtunnel` | The release tag to upload to |
| `mytunnel-mac mytunnel-mac-arm64 mytunnel-linux` | Files to upload |
| `--clobber` | Overwrite existing assets with the same name |

### Create a new release

```bash
gh release create v0.2.0 mytunnel-mac mytunnel-mac-arm64 mytunnel-linux \
  --title "v0.2.0" \
  --notes "Release notes here"
```

---

## Git Commands

### Check status

```bash
git status
```

Shows modified, staged, and untracked files.

### Stage and commit

```bash
git add .
git commit -m "your message"
```

### Push to GitHub

```bash
git push origin master
```

### View commit history

```bash
git log --oneline
```

---

## Install Script (one-liner)

```bash
curl -fsSL https://raw.githubusercontent.com/DpkRn/devtunnel/master/install.sh | bash
```

| Flag | Meaning |
|------|---------|
| `-f` | Fail silently on HTTP errors (non-zero exit) |
| `-s` | Silent mode — no progress bar |
| `-S` | Show error even in silent mode |
| `-L` | Follow redirects |

The script detects OS and CPU architecture automatically:
- `uname` → gets the OS (`Darwin` or `Linux`)
- `uname -m` → gets the CPU arch (`arm64` or `x86_64`)
- Downloads the correct binary and moves it to `/usr/local/bin/`

---

## Debugging a Bad Install

```bash
# Check what's actually installed
file /usr/local/bin/mytunnel

# Check CPU architecture of your Mac
uname -m

# Manually download and inspect before installing
curl -fsSL <URL> -o /tmp/test-binary
file /tmp/test-binary
xxd /tmp/test-binary | head -3
```
