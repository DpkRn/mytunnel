package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:9000")
	fmt.Println("Connected to tunnel server")

	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Disconnected from server")
			return
		}

		req := string(buffer[:n])
		fmt.Println("Received request:", req)

		parts := strings.Split(req, " ")
		// method := parts[0]
		path := parts[3]
		fmt.Println("Path:", path)
		// Call local server
		resp, err := http.Get("http://localhost:8080/" + path)
		if err != nil {
			conn.Write([]byte("Error calling local server"))
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		conn.Write(body)
	}
}
