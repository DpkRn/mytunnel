package main

import (
	"os"
	"os/signal"

	"github.com/DpkRn/devtunnel/internal/client"
)

func main() {
	port := os.Args[2]
	client.Start(port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
