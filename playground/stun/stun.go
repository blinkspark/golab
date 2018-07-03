package main

import (
	"fmt"
	"net"
	"time"

	"github.com/blinkspark/golab/util"
	"github.com/ccding/go-stun/stun"
)

func main() {
	nat, host, err := stun.NewClient().Discover()
	util.CheckErr(err)
	fmt.Println(nat)
	fmt.Println(host)
}

func client(addrstr string) {
	addr, err := net.ResolveUDPAddr("udp", addrstr)
	util.CheckErr(err)

	laddr, err := net.ResolveUDPAddr("udp", ":0")
	util.CheckErr(err)

	conn, err := net.DialUDP("udp", laddr, addr)
	util.CheckErr(err)
	defer conn.Close()

	for i := 0; i < 10; i++ {
		fmt.Println("sending to", addrstr)
		conn.Write([]byte("test a"))
		time.Sleep(time.Second)
	}
}
