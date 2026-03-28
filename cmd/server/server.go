package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	clientConn = make(map[string]net.Conn)
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
	subdomain := generateID()

	mu.Lock()
	clientConn[subdomain] = conn
	mu.Unlock()

	publicUrl := subdomain + ".localhost:3000"

	_, err := conn.Write([]byte(publicUrl + "\n"))
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
	fmt.Println("clientConn:", clientConn)
	fmt.Println("subdomain:", subdomain)
	conn, ok := clientConn[subdomain]
	fmt.Println("conn:", conn)
	fmt.Println("ok:", ok)
	mu.RUnlock()

	if !ok {
		http.Error(w, "Tunnel not found", 404)
		return
	}

	reqData := fmt.Sprintf("%s|%s", r.Method, r.URL.String())

	_, err := conn.Write([]byte(reqData))
	if err != nil {
		http.Error(w, "Tunnel write failed", 500)
		return
	}

	buffer := make([]byte, 4096)

	// ⏱ timeout (important)
	// conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	n, err := conn.Read(buffer)
	if err != nil {
		http.Error(w, "Tunnel timeout", 504)
		return
	}
	fmt.Println("Response:", string(buffer[:n]))

	_, err = w.Write(buffer[:n])
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
