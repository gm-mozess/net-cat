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

	var tab []net.Conn
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

		maxConnectMutex.Lock()
		cnxA++
		maxConnectMutex.Unlock()
		tab = append(tab, conn)

		fmt.Println(conn)
		
		if len(tab) > 1 {
			MessageWriter(conn, tab)
			MessageReader(conn, tab)
		}else{
			fmt.Fprintln(conn, "You need to be a peer for chating!")
			fmt.Fprint(conn, "["+timer+"]["+userName+"]:"+message)
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

func MessageWriter(conn net.Conn, tab []net.Conn) {

	writer := bufio.NewWriter(conn)

	for _, cnx := range tab {
		if cnx != conn {
			_, err := writer.WriteString("[" + timer + "][" + userName + "]:" + message)
			if err != io.EOF {
				CatchError(err)
			}
		}
	}

	writer.Flush()
}

func MessageReader(conn net.Conn, tab []net.Conn) {
	message = Reader(conn)
	MessageWriter(conn, tab)
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
