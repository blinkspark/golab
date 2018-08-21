package main

import (
	"bufio"
	"context"
	"github.com/blinkspark/golab/mylibp2p"
	"github.com/blinkspark/golab/mylibp2p/config"
	"github.com/blinkspark/golab/util"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	mjson "github.com/multiformats/go-multicodec/json"
	"log"
)

const (
	Protocol = "/PE/0.0.1"
)

type Test struct {
	Message string
}

func main() {
	h, err := mylibp2p.InitHost()
	util.CheckErr(err)

	c, err := config.ReadConfig()
	util.CheckErr(err)
	if len(c.Peers) == 0 {
		savePeers(c, h)
		log.Println("Node init complete, please add \"Peers\" to config.json")
		return
	}

	//add peer info of config to peer store
	for _, tpi := range c.Peers {
		if tpi.Tag == "Self" {
			continue
		}
		id, err := peer.IDB58Decode(tpi.ID)
		util.CheckErr(err)

		maddrs := make([]multiaddr.Multiaddr, 0)
		for _, a := range tpi.Addrs {
			maddrs = append(maddrs, multiaddr.StringCast(a))
		}
		h.Peerstore().AddAddrs(id, maddrs, peerstore.PermanentAddrTTL)
	}

	h.SetStreamHandler(Protocol, func(stream net.Stream) {
		log.Println("got a stream!")
		receivePeerInfos(stream, h, c)
		savePeers(c, h)
		sendPeerInfos(stream, h, c)
		stream.Close()
	})

	for _, p := range h.Peerstore().Peers() {
		stream, err := h.NewStream(context.Background(), p, Protocol)
		if err != nil {
			log.Println(err)
			continue
		}

		sendPeerInfos(stream, h, c)
		receivePeerInfos(stream, h, c)
		savePeers(c, h)

		stream.Close()
	}

	select {}
}

func receivePeerInfos(stream net.Stream, h host.Host, c *config.Config) {
	reader := bufio.NewReader(stream)
	dec := mjson.Multicodec(false).Decoder(reader)
	tpi := new(config.TaggedPeerInfo)
	for {
		err := dec.Decode(tpi)
		if err != nil {
			break
		}
		log.Println(*tpi)
		id, err := peer.IDB58Decode(tpi.ID)
		util.CheckErr(err)
		for _, a := range tpi.Addrs {
			h.Peerstore().AddAddr(id, multiaddr.StringCast(a), peerstore.PermanentAddrTTL)
		}
	}
}

func sendPeerInfos(stream net.Stream, h host.Host, c *config.Config) {
	writer := bufio.NewWriter(stream)
	enc := mjson.Multicodec(false).Encoder(writer)
	for _, tpi := range c.Peers {
		if tpi.Tag == "Self" {
			tpi.Tag = ""
			maddrs := h.Addrs()
			addrs := make([]string, 0)
			for _, ma := range maddrs {
				addrs = append(addrs, ma.String())
			}
			tpi.Addrs = addrs
		}
		enc.Encode(tpi)
	}
	writer.Flush()
}

func savePeers(c *config.Config, h host.Host) {
	peers := h.Peerstore().Peers()
	tpis := make([]config.TaggedPeerInfo, 0)
	for _, p := range peers {
		pi := h.Peerstore().PeerInfo(p)
		addrs := make([]string, 0)
		for _, a := range pi.Addrs {
			addrs = append(addrs, a.String())
		}

		tag := ""
		if p == h.ID() {
			tag = "Self"
		}
		tpi := config.TaggedPeerInfo{ID: pi.ID.Pretty(), Addrs: addrs, Tag: tag}
		tpis = append(tpis, tpi)
	}
	c.Peers = tpis
	config.SaveConfig(c)
}
