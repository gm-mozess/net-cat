package main

import (
	"fmt"
	"log"
	"netcat/server"
	"os"
)

func main() {

	args := len(os.Args)

	switch args {
	case 1:
		server.Port = "8989"
	case 2:
		server.Port = os.Args[1]
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
	}

	go func() {
			for err := range server.ChanError{
				if err != nil{
					log.Fatal(err)
				}
			}
	}()

	server.ServerTCP()
	
}
