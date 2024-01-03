package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var Port string
var maxConnect int

func ServerTCP() {
	listener, err := net.Listen("tcp", ":"+Port)
	CatchError(err)

	fmt.Println("Listening on the port :" + Port)

	defer listener.Close()

	for {
		if maxConnect < 11{
			
			conn, err := listener.Accept()
			CatchError(err)
			maxConnect++
			go IncommingConnections(conn)
		}


	}

}

func IncommingConnections(conn net.Conn) {
	//time := time.Now().Format("01-01-1889 13:45:45 GHL")
	conn.Write(WelcomeMessage())
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	for {
		data := Reader(conn)
		fmt.Println([]byte(data))
		for data == "" {
			conn.Write([]byte("[ENTER YOUR NAME]: "))
			data = Reader(conn)
		}

	}
}

func Reader(conn net.Conn) string {

	// Read data from the client
	netData, err := bufio.NewReader(conn).ReadString('\n')
	netData = strings.Trim(netData, "\n")

	if err != nil {
		log.Fatal("Error:", err)
	}
	return netData
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
