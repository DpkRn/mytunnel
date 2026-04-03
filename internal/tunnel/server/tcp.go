package server

import (
	"fmt"
	"log"
	"net"

	"github.com/DpkRn/devtunnel/internal/config"
	"github.com/DpkRn/devtunnel/internal/id"
	"github.com/hashicorp/yamux"
)

// TCPListener accepts tunnel client connections on the control port.
type TCPListener struct {
	Registry *Registry
	Addr     string
}

// NewTCPListener returns a listener configured with defaults from config package.
func NewTCPListener(reg *Registry) *TCPListener {
	return &TCPListener{
		Registry: reg,
		Addr:     config.ControlListenAddr,
	}
}

// ListenAndServe blocks forever accepting clients.
func (t *TCPListener) ListenAndServe() error {
	ln, err := net.Listen("tcp", t.Addr)
	if err != nil {
		return fmt.Errorf("listen tcp %s: %w", t.Addr, err)
	}
	log.Printf("tunnel control listening on %s", t.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go t.handleConn(conn)
	}
}

func (t *TCPListener) handleConn(conn net.Conn) {
	sub := id.Generate()
	publicLine := sub + config.PublicHostSuffix + "\n"

	if _, err := conn.Write([]byte(publicLine)); err != nil {
		conn.Close()
		return
	}

	session, err := yamux.Server(conn, nil)
	if err != nil {
		conn.Close()
		return
	}

	t.Registry.Add(sub, session)
	log.Printf("tunnel client connected: %s", sub)
}
