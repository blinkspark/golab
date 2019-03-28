package main

import (
	"flag"
	"fmt"
	"github.com/blinkspark/golab/util"
	"log"
	"net"
	"os"
	"os/signal"
)

var PORT = flag.Int("p", 22332, "PORT of the server")

func main() {
	// solve tcp address
	flag.Parse()
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *PORT))
	util.CheckErr(err)

	// create a goroutine to handle serverSocket socket
	serverSocket, err := net.ListenTCP("tcp", addr)
	util.CheckErr(err)
	log.Println("listening")
	go serverSocketHandler(serverSocket)
	util.CheckErr(err)

	// listen signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	s := <-sigChan
	log.Println(s)

	// cleanup
	err = serverSocket.Close()
	util.CheckErr(err)
}

func serverSocketHandler(server *net.TCPListener) {
	log.Println(server.Addr())
	for {
		if conn, err := server.Accept(); err == nil {
			log.Println("accepted")
			go clientHandler(conn)
		}
	}
}

func clientHandler(conn net.Conn) {
	go readFromClient(conn)
}

func readFromClient(conn net.Conn) {
	buffer := make([]byte, 4*1024)
	for {
		if n, err := conn.Read(buffer); err == nil && n > 0 {
			data := buffer[:n]
			log.Println(string(data))
		}
	}
}
