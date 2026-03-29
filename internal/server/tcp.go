package server

import (
	"fmt"
	"net"

	"github.com/DpkRn/devtunnel/internal/pkg"
	"github.com/hashicorp/yamux"
)

func StartTCP(reg *Registry) {
	listener, _ := net.Listen("tcp", ":9000")

	for {
		conn, _ := listener.Accept()
		go handleClient(conn, reg)
	}
}

func handleClient(conn net.Conn, reg *Registry) {
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

	reg.Add(subdomain, session)

	fmt.Println("Client connected:", subdomain)
}
