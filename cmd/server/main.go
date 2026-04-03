package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/DpkRn/devtunnel/internal/server"
)

func main() {
	reg := server.NewRegistry()
	tcp := server.NewTCP(reg)
	go tcp.StartTCP()
	handler := server.NewHandler(reg)

	http.HandleFunc("/", handler.HandleRequest())
	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatalf("Failed to listen on port 3000: %v", err)
		}
	}()
	fmt.Println("✅HTTP Server Listening on port 3000")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
