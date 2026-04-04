# devtunnel

Expose local ports to the internet using a simple CLI.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/DpkRn/devtunnel/master/install.sh | bash
```

This will download the latest `mytunnel` binary for your OS (Linux or macOS) and install it to `/usr/local/bin`.

## Usage

```bash
mytunnel http 3000
```

## Build from Source

```bash
go build -o mytunnel ./cmd/client
sudo mv mytunnel /usr/local/bin/
```

For a completely fresh build:

```bash
go clean -cache -modcache -i -r
go build -a -o mytunnel ./cmd/client
```

## Docker (tunnel server + nginx)

- **443** — HTTPS (nginx → mytunneld :3000, `X-Forwarded-Proto: https`).
- **80** — redirects to HTTPS; `/.well-known/acme-challenge/` is served for Let’s Encrypt.
- **9000** — tunnel control (mytunneld), not nginx.

Certs in `nginx/ssl/`: a cert for **`clickly.cv` + `www` only** does **not** cover **`*.clickly.cv`** tunnel URLs — Chrome shows **ERR_CERT_COMMON_NAME_INVALID**. Use **`scripts/gen-ssl-selfsigned.sh`** (includes `*.clickly.cv` in SAN; still “not secure” until you trust the CA) or a **Let’s Encrypt wildcard** `*.clickly.cv` via **DNS-01** (certbot DNS plugin for your DNS host; HTTP-01 cannot issue wildcard). Then copy PEMs into `nginx/ssl/` and `docker compose restart nginx`.

Production domain is set in `internal/config/config.go` (`PublicHostSuffix`, `PublicURLScheme`). The CLI dials **`clickly.cv:9000`** for the tunnel; override with `DEVTUNNEL_SERVER=localhost:9000` for local dev.

```bash
chmod +x scripts/docker-server.sh && ./scripts/docker-server.sh
```

The script uses **`docker compose`** (plugin) if available, otherwise **`docker-compose`**.

Ubuntu’s default apt **does not** ship `docker-compose-plugin`. Install the Compose v2 binary once:

```bash
chmod +x scripts/install-docker-compose.sh && ./scripts/install-docker-compose.sh
```

Or add [Docker’s official apt repo](https://docs.docker.com/engine/install/ubuntu/) and install `docker-compose-plugin` from there.

Plain Docker (no nginx):

```bash
docker build -t mytunneld .
docker run --rm -p 3000:3000 -p 9000:9000 mytunneld
```

The mytunneld image is **distroless** (~10 MB). `docker image prune -f` clears dangling build layers.
