package main

import (
	"context"
	"github.com/blinkspark/golab/mylibp2p"
	"github.com/blinkspark/golab/util"
	"github.com/libp2p/go-floodsub"
	"log"
	"time"
)

func main() {
	h, err := mylibp2p.InitHost()
	util.CheckErr(err)

	log.Println(h.ID().Pretty())
	log.Println(h.Addrs())

	ps, err := floodsub.NewFloodSub(context.Background(), h)
	util.CheckErr(err)

	for {
		err = ps.Publish("test", []byte("Hello"))
		util.CheckErr(err)
		log.Println("published")
		time.Sleep(time.Second)
	}
}
