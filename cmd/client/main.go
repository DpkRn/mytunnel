package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/DpkRn/devtunnel/internal/client"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: mytunnel http <port>")
		return
	}
	command := os.Args[1]
	port := os.Args[2]

	switch command {
	case "http":
		client.Start(port)
	default:
		fmt.Println("Unknown command:", command)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
