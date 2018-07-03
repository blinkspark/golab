package main

import (
	"fmt"
	"net"
	"time"

	"github.com/blinkspark/golab/util"
)

func client() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:62233")
	util.CheckErr(err)

	laddr, err := net.ResolveUDPAddr("udp", ":0")
	util.CheckErr(err)

	conn, err := net.DialUDP("udp", laddr, addr)
	util.CheckErr(err)
	defer conn.Close()

	for i := 0; i < 10; i++ {
		conn.Write([]byte("test a"))
		time.Sleep(time.Second)
	}
}

func server() {
	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:62233")
	util.CheckErr(err)

	conn, err := net.ListenUDP("udp", laddr)
	util.CheckErr(err)

	for {
		buffer := make([]byte, 4096)
		n, c, err := conn.ReadFromUDP(buffer)
		util.CheckErr(err)

		fmt.Println(string(buffer[0:n]), "form", c.String())
	}
}

func main() {
	go client()
	server()
}
