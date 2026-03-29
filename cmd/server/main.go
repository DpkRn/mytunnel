package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/DpkRn/devtunnel/internal/server"
)

func main() {
	reg := server.NewRegistry()

	go server.StartTCP(reg)

	http.HandleFunc("/", server.Handler(reg))
	http.ListenAndServe(":3000", nil)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
