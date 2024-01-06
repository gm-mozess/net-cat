package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Port            string
	maxConnectMutex sync.Mutex
	userName        string
	message         string
	cnxA            int
	clients         []net.Conn
)

var timer = time.Now().Format("2006-01-02 15:04:05")

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

	if cnxA > 10 {
		return
	} else {
		fmt.Fprint(conn, string(WelcomeMessage()))
		userName = ""
		//conn.Write(WelcomeMessage())
		for userName == "" {
			fmt.Fprint(conn, "[ENTER YOUR NAME]: ")
			userName = Reader(conn)
		}

		if cnxA > 1 {
			fmt.Fprint(conn, "["+timer+"]["+userName+"]:")
		} else {
			fmt.Fprintln(conn, "You need to be a peer for chating!")
			fmt.Fprint(conn, "["+timer+"]["+userName+"]:")
		}

	  	maxConnectMutex.Lock()
		cnxA++
		clients = append(clients, conn)
		maxConnectMutex.Unlock()
		MessageWriter(conn)
	}
}

func Reader(conn net.Conn) string {

	// Read data from the client
	netData, err := bufio.NewReader(conn).ReadString('\n')
	netData = strings.Trim(netData, "\n")

	if err != nil {
		if err == io.EOF {
			return "/logout"
		} else {
			log.Fatal("Error:", err)
		}
	}
	return netData
}

func MessageWriter(conn net.Conn) {
	message = Reader(conn)
	BroadcastMessage(message, conn)

}

func BroadcastMessage(msg string, sender net.Conn) {
	fmt.Println(len(clients))
	// Iterate over all connected clients and send the message
	for _, client := range clients {
		if client != sender {
			writer := bufio.NewWriter(client)
			_, err := writer.WriteString("["+timer+"]["+userName+"]:"+msg+"\n")
			CatchError(err)
			writer.Flush()
		}
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
