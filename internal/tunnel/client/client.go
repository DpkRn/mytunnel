package client

import (
	"fmt"
	"net"
	"strings"

	"github.com/DpkRn/devtunnel/internal/config"
	"github.com/hashicorp/yamux"
)

// Options configure the tunnel client.
type Options struct {
	// ServerAddr is the tunnel control plane (yamux), e.g. "host:9000".
	ServerAddr string
	// LocalHost is the host used when forwarding to the local HTTP server (default "localhost").
	LocalHost string
}

// Connect dials the tunnel server, starts stream forwarding to localhost:<port>, and returns
// the public URL plus a stop function. Non-blocking except for the initial handshake.
func Connect(port string, opts Options) (publicURL string, stop func(), err error) {
	if opts.ServerAddr == "" {
		opts.ServerAddr = config.DefaultTunnelServerAddr
	}
	if opts.LocalHost == "" {
		opts.LocalHost = "localhost"
	}

	conn, err := net.Dial("tcp", opts.ServerAddr)
	if err != nil {
		return "", noop, fmt.Errorf("dial tunnel server %s: %w", opts.ServerAddr, err)
	}

	session, err := yamux.Client(conn, nil)
	if err != nil {
		conn.Close()
		return "", noop, fmt.Errorf("yamux client: %w", err)
	}

	go acceptStreams(session, opts.LocalHost, port)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		session.Close()
		conn.Close()
		return "", noop, fmt.Errorf("read public URL: %w", err)
	}

	url := "http://" + strings.TrimSpace(string(buf[:n]))
	stop = func() {
		session.Close()
		conn.Close()
	}
	return url, stop, nil
}

func noop() {}
