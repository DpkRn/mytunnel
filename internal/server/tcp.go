package server

import (
	"fmt"
	"log"
	"net"

	"github.com/DpkRn/devtunnel/internal/pkg"
	"github.com/hashicorp/yamux"
)

type TCP interface {
	StartTCP() error
	HandleClient(conn net.Conn, reg *Registry)
}

type tcp struct {
	Registry *Registry
}

func NewTCP(reg *Registry) TCP {
	return &tcp{
		Registry: reg,
	}
}

func (t *tcp) StartTCP() error {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
		return err
	}
	fmt.Println("✅TCP Connection Listening on port 9000")

	for {
		conn, _ := listener.Accept()
		go t.HandleClient(conn, t.Registry)
	}
}

func (t *tcp) HandleClient(conn net.Conn, reg *Registry) {
	subdomain := pkg.GenerateID()
	publicUrl := subdomain + ".localhost:3000\n"

	// ✅ send BEFORE yamux
	conn.Write([]byte(publicUrl))

	// now start yamux
	session, err := yamux.Server(conn, nil)
	if err != nil {
		conn.Close()
		return
	}

	t.Registry.Add(subdomain, session)

	fmt.Println("Client connected:", subdomain)
}
