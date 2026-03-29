package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {

	fmt.Println(os.Args[0], os.Args[1])
	if len(os.Args) < 3 {
		fmt.Println(os.Args)
		fmt.Println("Usage: mytunnel http <port>")
		return
	}
	protocol := os.Args[1]
	port := os.Args[2]
	// protocol := "http"
	// port := "8080"
	startTunneling(protocol, port)
}
func startTunneling(protocol, port string) {
	if protocol != "http" {
		fmt.Println("only http protocol supported")
		return
	}
	fmt.Println("Connecting to server...")
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	//read generatedsubdomain
	publicUrl := make([]byte, 4096)
	n, err := conn.Read(publicUrl)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	publicUrl = []byte(strings.TrimSpace(string(publicUrl[:n])))

	fmt.Println("🌐 Public URL → http://" + string(publicUrl))

	//make request on local server
	localServer := "http://localhost:" + port

	for {
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		request := buffer[:n]
		// parts := strings.Split(request, "|")
		// path := parts[1]
		// fmt.Println("path:", path)
		tunnelRequest := TunnelRequest{}
		err = json.Unmarshal(request, &tunnelRequest)
		if err != nil {
			fmt.Println("Error unmarshalling request:", err)
			return
		}
		fmt.Println("tunnelRequest:", tunnelRequest)
		fmt.Println("body:", string(tunnelRequest.Body))
		reqUrl, err := http.NewRequest(tunnelRequest.Method, localServer+tunnelRequest.Path, bytes.NewBuffer(tunnelRequest.Body))
		reqUrl.Header = tunnelRequest.Headers
		client := &http.Client{}
		resp, err := client.Do(reqUrl)
		if err != nil {
			fmt.Println("Error making request to local server:", err)
			return
		}
		defer resp.Body.Close()
		respObj := TunnelResponse{
			Status:  resp.StatusCode,
			Headers: resp.Header,
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response from local server:", err)
			return
		}
		respObj.Body = body
		fmt.Println("Response:", string(body))

		responseData, err := json.Marshal(respObj)
		conn.Write(responseData)
	}
}

type TunnelRequest struct {
	Method  string
	Path    string
	Headers http.Header
	Body    []byte
}

type TunnelResponse struct {
	Status  int
	Headers http.Header
	Body    []byte
}
