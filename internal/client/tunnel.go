package client

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/yamux"
)

func Start(port string) {
	conn, _ := net.Dial("tcp", "localhost:9000")

	session, _ := yamux.Client(conn, nil)

	go acceptStreams(session, port)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	fmt.Println("Public URL:", strings.TrimSpace(string(buf[:n])))
}
