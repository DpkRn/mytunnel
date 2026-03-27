// client.go
package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to server")

	// Read messages from server
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Server closed connection")
				return
			}
			fmt.Println("Server says:", string(buffer[:n]))
		}
	}()

	// Send messages to server
	for {
		var input string
		fmt.Scanln(&input)
		fmt.Println("Enter message to send to server")
		conn.Write([]byte(input))
	}
}
