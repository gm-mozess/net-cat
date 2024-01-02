package server

import (
	"fmt"
	"net"
	"os"
)

var ChanError = make(chan error)
var Port string

func ServerTCP() {
	listener, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		ChanError <- err
		
	}

	fmt.Println("Listening on the port :" + Port)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			ChanError <- err
			return
			
		}
		go IncommingConnections(conn)
	}
}

func IncommingConnections(conn net.Conn) {
	//time := time.Now().Format("01-01-1889 13:45:45 GHL")
	conn.Write(WelcomeMessage())

}

func WelcomeMessage() []byte {
	file, err := os.Open("./pingoin.txt")
	ChanError <- err

	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)

	if err != nil {
		ChanError <- err
		return nil
	}
	return buffer[:n]
}
