package main

import (
	"context"
	"github.com/blinkspark/golab/util"
	"github.com/ipfs/go-ipfs-addr"
	"github.com/libp2p/go-floodsub"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-peerstore"
	"log"
	"time"
)

func main() {
	h, err := libp2p.New(context.Background())
	util.CheckErr(err)

	ps, err := floodsub.NewFloodSub(context.Background(), h)
	util.CheckErr(err)

	addr, err := ipfsaddr.ParseString("/ip4/127.0.0.1/tcp/22330/ipfs/QmXGmsUea2HLj6bbLD2eWcWLzocxe9iZBtVhPXHT4VAbMb")
	util.CheckErr(err)

	pi, err := peerstore.InfoFromP2pAddr(addr.Multiaddr())
	if err := h.Connect(context.Background(), *pi); err != nil {
		log.Println(err)
	} else {
		log.Println("Connection established with bootstrap node: ", *pi)
	}

	s, err := ps.Subscribe("test")
	util.CheckErr(err)

	for {
		func() {
			t10sctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
			defer cancel()
			m, err := s.Next(t10sctx)
			util.CheckErr(err)

			log.Println("From: ", m.GetFrom())
			log.Println("Data: ", string(m.GetData()))
			log.Println("seq: ", m.GetSeqno())
			log.Println("topic: ", m.GetTopicIDs())
		}()
	}

}
