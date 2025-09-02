package main

import (
	"fmt"
	"netcat/server"
	"os"
)

func main() {
	args := len(os.Args)
	var port string

	switch args {
	case 1:
		port = "8989" // Default port
	case 2:
		port = os.Args[1]
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	nc := server.NewChatServer(port)
	// Start the TCP server
	nc.Start()
}
