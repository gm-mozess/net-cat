package server

import (
	"net"
)

var GlobalErr = make(chan error)

func ServerHandler() {
	listener, err := net.Listen("tcp", "localhost:8080")
	GlobalErr <- err

	defer listener.Close()

	for {

		conn, err := listener.Accept()
		GlobalErr <- err

		go ClientHandler(conn)
	}
}

func ClientHandler(conn net.Conn) {
	defer conn.Close()

}
