package server

import (
	"fmt"
	"log"
	"net"
	"os"
)

var Port string

func ServerTCP() {
	listener, err := net.Listen("tcp", ":"+Port)
	CatchError(err)

	fmt.Println("Listening on the port :" + Port)

	defer listener.Close()

	for {

		conn, err := listener.Accept()
		CatchError(err)
		go IncommingConnections(conn)
	}

}

func IncommingConnections(conn net.Conn) {
	//time := time.Now().Format("01-01-1889 13:45:45 GHL")
	conn.Write(WelcomeMessage())

}

func WelcomeMessage() []byte {
	file, err := os.Open("./pingoin.txt")
	CatchError(err)

	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	CatchError(err)

	return buffer[:n]
}

func CatchError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
