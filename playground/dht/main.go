package main

import (
	"bufio"
	"context"
	"github.com/blinkspark/golab/util"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs-addr"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multihash"
	"log"
	"time"
)

var bootstrapPeers = []string{
	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
}

const TestProtocol = "/TestDHT/0.0.1"

func main() {
	h, err := libp2p.New(context.Background())
	util.CheckErr(err)

	h.SetStreamHandler(TestProtocol, handleStream)

	dht := dht.NewDHT(context.Background(), h, datastore.NewMapDatastore())
	util.Ignore(dht)

	for _, p := range bootstrapPeers {
		addr, err := ipfsaddr.ParseString(p)
		util.CheckErr(err)

		pi, err := peerstore.InfoFromP2pAddr(addr.Multiaddr())
		if err := h.Connect(context.Background(), *pi); err != nil {
			log.Println(err)
		} else {
			log.Println("Connection established with bootstrap node: ", *pi)
		}

	}
	cidBuilder := cid.V1Builder{Codec: cid.Raw, MhType: multihash.SHA2_256, MhLength: 0}
	cid, err := cidBuilder.Sum([]byte(TestProtocol))

	tctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = dht.Provide(tctx, cid, true)
	util.CheckErr(err)

	tctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	peers, err := dht.FindProviders(tctx, cid)

	for _, p := range peers {
		log.Println(p)
		if p.ID == h.ID() || len(p.Addrs) == 0 {
			continue
		}
		t10sctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		log.Println("Connecting to:", p)
		err := h.Connect(t10sctx, p)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Opening stream:")
		s, err := h.NewStream(context.Background(), p.ID, TestProtocol)
		util.CheckErr(err)
		sendMsg(s)
		receiveMsg(s)
		s.Close()
	}
	select {}
}

func handleStream(stream net.Stream) {
	log.Println(stream)
	receiveMsg(stream)
	sendMsg(stream)
}

func receiveMsg(stream net.Stream) {
	reader := bufio.NewReader(stream)
	buf := make([]byte, 1024)
	for {

		n, err := reader.Read(buf)
		if err != nil {
			log.Println("read stream terminated:", err)
			break
		}
		log.Println(string(buf[:n]))
	}
}

func sendMsg(stream net.Stream) {
	writer := bufio.NewWriter(stream)
	writer.WriteString("Hello\n")
	writer.Flush()
}
