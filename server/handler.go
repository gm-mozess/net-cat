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
)

var (
	Port string
	maxConnect int
	maxConnectMutex sync.Mutex
)

func ServerTCP() {
	listener, err := net.Listen("tcp", ":"+Port)
	CatchError(err)

	fmt.Println("Listening on the port :" + Port)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		CatchError(err)
		go IncommingConnections(conn)
		go Writer(conn)
	}
}

func IncommingConnections(conn net.Conn) {
	//time := time.Now().Format("01-01-1889 13:45:45 GHL")
	conn.Write(WelcomeMessage())
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	userName := Reader(conn)

	for userName == "" {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		userName = Reader(conn)
	}

	for {
		if maxConnect < 11 {

			if userName == "/logout" {
				maxConnectMutex.Lock()
				maxConnect--
				maxConnectMutex.Unlock()
				return
			}
			
			maxConnectMutex.Lock()
			maxConnect++
			maxConnectMutex.Unlock()

			conn.Write([]byte("["+userName+"] Enter a message (/logout to exit): "))
			
		}else{
			conn.Write([]byte("Maximum connection reached!"))
			return
		}
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


func Writer(conn net.Conn) {

	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(conn)

	for scanner.Scan(){
		message := scanner.Text()

		_,err := writer.WriteString(message+"\n")
		if err != io.EOF{
			CatchError(err)
		}
	}
	writer.Flush()
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
