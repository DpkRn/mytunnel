# Tunnel client (`server.go`)

This folder holds the **client** that connects to your tunnel server, forwards HTTP from the tunnel to a **local** app, and sends responses back. You run it on the machine where your app listens (e.g. `localhost:3000`).

---

## What you pass on the command line

| Argument | Meaning |
|----------|---------|
| `http` | Protocol (only `http` is supported in code). |
| `<port>` | Port your **local** server uses (e.g. `3000` → traffic goes to `http://localhost:3000`). |

Replace `<port>` with a real number, e.g. `3000`.

---

## Run without installing (development)

From this `client` directory:

```bash
go run server.go http 3000
```

**Why three “words” after `go run`?** Go passes them as `os.Args`:

- `os.Args[0]` → `server.go` (the file being run)
- `os.Args[1]` → `http`
- `os.Args[2]` → `3000`

So the program expects **at least** two user arguments after the program name: protocol + port.

---

## Build a binary you can run by name

```bash
go build -o mytunnel server.go
```

- **`go build`** compiles `server.go` into an executable.
- **`-o mytunnel`** names that executable `mytunnel` (instead of `server` or `client`).

Run it from the same directory:

```bash
./mytunnel http 3000
```

Here `os.Args[0]` is `./mytunnel` (the binary path), not `server.go`.

---

## Run `mytunnel` from anywhere (optional, macOS/Linux)

Put the binary on your **PATH** so the shell finds it without `./`:

```bash
sudo mv mytunnel /usr/local/bin/
```

Then from any directory:

```bash
mytunnel http 3000
```

`/usr/local/bin` is a common place for user-installed tools; `sudo` is needed because that directory is usually owned by root.

---

## Quick recap

1. **`go run server.go http <port>`** — quick test, no install.
2. **`go build -o mytunnel server.go`** then **`./mytunnel http <port>`** — same behavior, reusable binary.
3. **`sudo mv mytunnel /usr/local/bin/`** — call **`mytunnel http <port>`** from any folder.
