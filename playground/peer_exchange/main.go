package main

import (
	"context"
	"encoding/json"
	"github.com/blinkspark/golab/mylibp2p"
	"github.com/blinkspark/golab/mylibp2p/config"
	"github.com/blinkspark/golab/util"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peerstore"
	"log"
)

const (
	Protocol = "/PE/0.0.1"
)

func main() {
	h, err := mylibp2p.InitHost()
	util.CheckErr(err)

	c, err := config.ReadConfig()
	util.CheckErr(err)
	if len(c.Peers) == 0 {
		initConfig(c, h)
		log.Println("Node init complete, please add \"Peers\" to config.json")
		return
	}

	//add peer info of config to peer store
	for _, pi := range c.Peers {
		if pi.Tag == "Self" {
			continue
		}
		h.Peerstore().AddAddrs(pi.PeerInfo.ID, pi.PeerInfo.Addrs, peerstore.PermanentAddrTTL)
	}

	h.SetStreamHandler(Protocol, func(stream net.Stream) {
		log.Println("got a stream!")
		buf := make([]byte, 1024)
		for i := 0; i < 3; i++ {
			n, err := stream.Read(buf)
			if err != nil {
				break
			}
			log.Println(string(buf[:n]))
		}
		stream.Close()
	})

	for _, p := range h.Peerstore().Peers() {
		stream, err := h.NewStream(context.Background(), p, Protocol)
		if err != nil {
			continue
		}

		for _, innerP := range h.Peerstore().Peers() {
			pi := h.Peerstore().PeerInfo(innerP)
			data, err := json.Marshal(pi)
			if err != nil {
				continue
			}
			stream.Write(data)
		}

		stream.Close()
	}

	select {}
}

func initConfig(c *config.Config, h host.Host) {
	peers := h.Peerstore().Peers()
	for _, p := range peers {
		pi := h.Peerstore().PeerInfo(p)
		tpi := config.TaggedPeerInfo{PeerInfo: pi, Tag: "Self"}
		c.Peers = append(c.Peers, tpi)
	}
	config.SaveConfig(c)
}
