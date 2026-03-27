package main

import (
	"fmt"
	"net"
)

func main() {

	//for tcp
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on :9000")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// defer conn.Close()

	//read from connection
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	fmt.Println("Received data:", string(buffer))
	conn.Write([]byte("Message received"))
	conn.Write([]byte("Message received part 2"))
}
