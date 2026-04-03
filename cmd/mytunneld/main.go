package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/DpkRn/devtunnel/internal/config"
	tunnelserver "github.com/DpkRn/devtunnel/internal/tunnel/server"
)

func main() {
	reg := tunnelserver.NewRegistry()

	tcp := tunnelserver.NewTCPListener(reg)
	go func() {
		if err := tcp.ListenAndServe(); err != nil {
			log.Fatalf("control plane: %v", err)
		}
	}()

	edge := tunnelserver.NewEdgeHTTP(reg)
	http.HandleFunc("/", edge.Handler())

	go func() {
		log.Printf("edge HTTP listening on %s", config.HTTPListenAddr)
		if err := http.ListenAndServe(config.HTTPListenAddr, nil); err != nil {
			log.Fatalf("edge HTTP: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
