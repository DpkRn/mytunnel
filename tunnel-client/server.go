package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/yamux"
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
	session, err := yamux.Client(conn, nil)
	if err != nil {
		fmt.Println("Error creating session:", err)
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

	for {
		stream, err := session.Accept()
		if err != nil {
			fmt.Println("Error opening stream:", err)
			return
		}
		go handleStream(stream, port)

	}
}

func handleStream(stream net.Conn, port string) {

	defer stream.Close()
	localServer := "http://localhost:" + port

	reader := bufio.NewReader(stream)
	requestByte, err := reader.ReadBytes('\n')

	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	fmt.Println("requestByte:", string(requestByte))
	tunnelRequest := TunnelRequest{}
	err = json.Unmarshal(requestByte, &tunnelRequest)
	if err != nil {
		fmt.Println("Error unmarshalling request:", err)
		return
	}
	fmt.Println("tunnelRequest:", tunnelRequest)
	fmt.Println("body:", string(tunnelRequest.Body))
	reqUrl, err := http.NewRequest(tunnelRequest.Method, localServer+tunnelRequest.Path, bytes.NewBuffer(tunnelRequest.Body))

	for k, v := range tunnelRequest.Headers {
		for _, val := range v {
			reqUrl.Header.Add(k, val)
		}
	}

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
	stream.Write(append(responseData, '\n'))
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
