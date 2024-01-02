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
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	buffer := make([]byte, 1024)
	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
		 fmt.Println("Error:", err)
		 return
		}
	  
		// Process and use the data (here, we'll just print it)
		fmt.Printf("Received: %s", buffer[:n])
	}
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
