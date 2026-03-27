package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println(os.Args)
		fmt.Println("Usage: mytunnel http <port>")
		return
	}
	protocol := os.Args[1]
	port := os.Args[2]
	startTunneling(protocol, port)
}
func startTunneling(protocol, port string) {
	if protocol != "http" {
		fmt.Println("only http protocol supported")
		return
	}
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	fmt.Println("✅ Connected to tunnel server")
	fmt.Println("🚀 Forwarding → http://localhost:", port)
	fmt.Println("🌐 Public URL → http://localhost:3000")
	defer conn.Close()
	for {
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		request := string(buffer[:n])
		fmt.Println("Request:", request)
		parts := strings.Split(request, " ")
		// method := parts[0]
		path := parts[1]

		//make request on local server
		localServer := "http://localhost:" + port

		resp, err := http.Get(localServer + path)
		if err != nil {
			fmt.Println("Error making request to local server:", err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response from local server:", err)
			return
		}
		conn.Write([]byte(body))
	}
}
