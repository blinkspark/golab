package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/blinkspark/golab/util"
	"golang.org/x/net/ipv4"
)

func client(ifi *net.Interface) {
	raddr, err := net.ResolveUDPAddr("udp", "224.0.1.1:62233")
	util.CheckErr(err)

	laddr, err := net.ResolveUDPAddr("udp", ":12233")
	util.CheckErr(err)

	conn, err := net.DialUDP("udp", laddr, raddr)
	util.CheckErr(err)
	defer conn.Close()

	for i := 0; i < 10; i++ {
		conn.Write([]byte("test a"))
		time.Sleep(time.Second)
	}
}

func server(ifi *net.Interface, name string) {
	laddr, err := net.ResolveUDPAddr("udp", "224.0.1.1:62233")
	util.CheckErr(err)

	// conn, err := net.ListenUDP("udp", laddr)

	conn, err := net.ListenMulticastUDP("udp", nil, laddr)
	util.CheckErr(err)
	defer conn.Close()

	for {
		buffer := make([]byte, 4096)
		n, c, err := conn.ReadFromUDP(buffer)
		util.CheckErr(err)

		fmt.Println(name, string(buffer[0:n]), "form", c.String())
	}
}

func main() {
	ifi := flag.String("ifi", "", "-ifi=en2")
	flag.Parse()
	showInterfaces()
	en2, err := net.InterfaceByName(*ifi)
	util.CheckErr(err)
	go client(en2)
	server(en2, "s1")
}

func showInterfaces() {
	ifis, err := net.Interfaces()
	util.CheckErr(err)
	for _, i := range ifis {
		fmt.Println(i)
		addrs, err := i.Addrs()
		util.CheckErr(err)
		for _, a := range addrs {
			fmt.Println("	", a.Network(), a)
		}
	}
}

func udpMulticastSender(groupAddr string, ifi *net.Interface) {
	conn, err := net.ListenPacket("udp", ":0")
	util.CheckErr(err)
	defer conn.Close()

	gipaddr, err := net.ResolveUDPAddr("udp", groupAddr)
	util.CheckErr(err)
	gip := gipaddr.IP
	fmt.Println(gip)

	pconn := ipv4.NewPacketConn(conn)
	en2, err := net.InterfaceByName("en2")
	util.CheckErr(err)

	err = pconn.JoinGroup(en2, gipaddr)
	if err != nil {
		err = pconn.SetControlMessage(ipv4.FlagDst, true)
		util.CheckErr(err)

		buffer := make([]byte, 1500)
		for {
			n, cm, src, err := pconn.ReadFrom(buffer)
			util.CheckErr(err)
			fmt.Println(n, cm, src)
			fmt.Println(buffer[:n])
		}
	}

}
