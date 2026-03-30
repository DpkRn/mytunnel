package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/DpkRn/devtunnel/internal/client"
)

func printHelp() {
	fmt.Println("mytunnel — expose your local server to the internet")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  mytunnel <command> <port>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  http <port>   Forward HTTP traffic to localhost:<port>")
	fmt.Println("  help          Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mytunnel http 3000")
	fmt.Println("  mytunnel http 8080")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	if command == "help" || command == "--help" || command == "-h" {
		printHelp()
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: mytunnel http <port>")
		return
	}

	port := os.Args[2]

	switch command {
	case "http":
		client.Start(port)
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Run 'mytunnel help' to see available commands.")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
