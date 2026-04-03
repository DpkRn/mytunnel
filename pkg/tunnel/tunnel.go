// Package tunnel is the stable public API for embedding devtunnel in other Go programs.
package tunnel

import (
	tunnelclient "github.com/DpkRn/devtunnel/internal/tunnel/client"
)

// Start connects to the default tunnel server and forwards HTTP to localhost:<port>.
// Returns the public URL, a stop function, and an error. Fully non-blocking after connect.
func Start(port string, serverAddr ...string) (url string, stop func(), err error) {
	opts := tunnelclient.Options{}
	if len(serverAddr) > 0 && serverAddr[0] != "" {
		opts.ServerAddr = serverAddr[0]
	}
	return tunnelclient.Connect(port, opts)
}
