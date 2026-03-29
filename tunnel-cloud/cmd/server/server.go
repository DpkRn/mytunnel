package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/hashicorp/yamux"
)

var (
	clientConn = make(map[string]*yamux.Session)
	mu         sync.RWMutex
)

func main() {

	go startTcpServer()
	http.HandleFunc("/", handleHttp)
	fmt.Println("Server listening on :3000")
	http.ListenAndServe(":3000", nil)
}

func startTcpServer() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("TCP Server listening on :9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {

	session, err := yamux.Server(conn, nil)
	if err != nil {
		conn.Close()
		return
	}

	subdomain := generateID()

	mu.Lock()
	clientConn[subdomain] = session
	mu.Unlock()

	publicUrl := subdomain + ".localhost:3000"
	_, err = conn.Write([]byte(publicUrl + "\n"))
	if err != nil {
		conn.Close()
		return
	}

	fmt.Println("Client connected:", subdomain)
	// ❌ NO READ LOOP HERE
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(200)
		return
	}
	host := r.Host
	subdomain := strings.Split(host, ".")[0]
	fmt.Println("host:", host)
	mu.RLock()
	session, ok := clientConn[subdomain]
	mu.RUnlock()

	if !ok {
		http.Error(w, "Tunnel not found", 404)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", 500)
		return
	}
	tunnelRequest := TunnelRequest{
		Method:  r.Method,
		Path:    r.URL.String(),
		Headers: r.Header,
		Body:    body,
	}

	fmt.Println("body:", string(body))

	requestData, err := json.Marshal(tunnelRequest)
	if err != nil {
		http.Error(w, "Error marshalling request", 500)
		return
	}
	// reqData := fmt.Sprintf("%s", requestData)
	fmt.Println("requestData:", string(requestData))

	stream, err := session.Open()
	if err != nil {
		http.Error(w, "Stream error", 500)
		return
	}
	defer stream.Close()
	_, err = stream.Write([]byte(append(requestData, '\n')))
	if err != nil {
		http.Error(w, "Tunnel write failed", 500)
		return
	}

	reader := bufio.NewReader(stream)
	responseByte, err := reader.ReadBytes('\n')

	if err != nil {
		http.Error(w, "Tunnel timeout", 504)
		return
	}
	//response
	respObj := TunnelResponse{}
	if err := json.Unmarshal(responseByte, &respObj); err != nil {
		fmt.Println("Error while unmarshaling response")
	}
	for k, v := range respObj.Headers {
		for _, val := range v {
			w.Header().Add(k, val)
		}
	}

	w.WriteHeader(respObj.Status)
	_, err = w.Write(respObj.Body)
	if err != nil {
		fmt.Println("Write error:", err)
	}
}

func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	length := 6 + rand.Intn(3)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
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
