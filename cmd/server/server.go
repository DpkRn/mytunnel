package main

import (
	"fmt"
	"net"
	"net/http"
)

var clientConn net.Conn

func main() {

	//for tcp
	go startTcpServer()
	http.HandleFunc("/ping", handleHttp)
	fmt.Println("Server listening on :3000")
	http.ListenAndServe(":3000", nil)
}
func startTcpServer() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println(" Tcp Server listening on :9000")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		clientConn = conn
	}
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	// defer conn.Close()
	if clientConn == nil {
		w.Write([]byte("no connections"))
	}
	reqData := fmt.Sprintf("method: %s, url: %s", r.Method, r.URL.String())
	clientConn.Write([]byte(reqData))

	buffer := make([]byte, 4096)
	n, _ := clientConn.Read(buffer)

	w.Write(buffer[:n])
}
