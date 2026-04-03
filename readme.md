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

Run the server (`mytunneld`) on the machine that terminates public traffic before clients can connect.

## Use as a Go library

```go
import "github.com/DpkRn/devtunnel/pkg/tunnel"

url, stop, err := tunnel.Start("3000")
if err != nil {
    log.Fatal(err)
}
defer stop()
// url is the public http://… address; your app keeps running until you call stop()
```

## Build from Source

```bash
go build -a -o mytunnel ./cmd/mytunnel
go build -a -o mytunneld ./cmd/mytunneld
sudo mv mytunnel /usr/local/bin/
```

For a completely fresh build:

```bash
go clean -cache -modcache -i -r
go build -a -o mytunnel ./cmd/mytunnel
```
