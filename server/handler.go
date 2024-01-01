package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var ChanError = make(chan error)
var GlobalErr error
var Port string

func ServerTcp() {
	fmt.Println("Listening on the port :"+Port)
	listener, err := net.Listen("tcp", ":"+Port)
	GlobalErr = err
	

	defer listener.Close()

	for {

		conn, err := listener.Accept()
		GlobalErr = err

		go SendToServer(conn)
	}
}

/* func ClientTcp() {
	conn, err := net.Dial("tcp", ":"+Port)
	GlobalErr <- err

	defer conn.Close()

	SendToServer(conn)
	ReadFromServer(conn)
} */


func SendToServer(conn net.Conn) {

	scanner := bufio.NewScanner(os.Stdin)
	_ , err := conn.Write(scanner.Bytes())
	GlobalErr = err

	fmt.Println(string(scanner.Bytes()))
	
}
