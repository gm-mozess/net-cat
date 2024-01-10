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
	cnxA            int
)

var clients = make(map[net.Conn]string)
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
		fmt.Fprintln(conn, string(WelcomeMessage()))
		userName = ""
		//conn.Write(WelcomeMessage())
		for userName == "" {
			_, err := fmt.Fprint(conn, "[ENTER YOUR NAME]: ")
			CatchError(err)
			userName = Reader(conn)
		}

		maxConnectMutex.Lock()
		cnxA++
		clients[conn] = userName
		maxConnectMutex.Unlock()
		LogSignal(conn)
		go BroadcastMessage(conn)
	}
}

func Reader(conn net.Conn) string {

	// Read data from the client
	netData, err := bufio.NewReader(conn).ReadString('\n')
	netData = strings.Trim(netData, "\n")

	if err != nil {
		if err == io.EOF {
			_, err := fmt.Fprintln(conn, "\n"+userName+" has left our chat...")
			CatchError(err)
		} else {
			log.Fatal("Error:", err)
		}
	}
	return netData
}

func MessageWriter(conn net.Conn) string {
	var msg string
	for msg == "" {
		_, err := fmt.Fprint(conn, "["+timer+"]["+clients[conn]+"]:")
		CatchError(err)
		msg = Reader(conn)
	}
	return msg
}


func BroadcastMessage(sender net.Conn) {
	for {
		var msg = MessageWriter(sender)
		// Iterate over all connected clients and send the message
		for conn := range clients {
			if conn != sender {
				_, err := fmt.Fprintln(conn, "\n["+timer+"]["+clients[sender]+"]:"+msg)
				CatchError(err)
			}
		}

		for conn := range clients {
			if conn != sender {
				_, err := fmt.Fprint(conn, "["+timer+"]["+clients[conn]+"]:")
				CatchError(err)
			}
		}
	}
}


func LogSignal(loger net.Conn) {
	for conn := range clients {
		if conn != loger {
			_, err := fmt.Fprintln(conn, "\n"+clients[loger]+" has joinded our chat...")
			CatchError(err)
			_, err = fmt.Fprint(conn, "["+timer+"]["+clients[conn]+"]:")
			CatchError(err)
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
