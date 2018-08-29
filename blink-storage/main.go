package main

import (
	"context"
	"log"
	"time"

	"github.com/blinkspark/golab/mylibp2p"
	"github.com/blinkspark/golab/util"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs-addr"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multihash"
)

const Protocol = "/TestDHT/0.0.2"

var bootstrapPeers = []string{
	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
}

func main() {
	h, err := mylibp2p.InitHost()
	util.CheckErr(err)

	h.SetStreamHandler(Protocol, streamHandler)
	dhtNode := dht.NewDHT(context.Background(), h, datastore.NewMapDatastore())

	bootstrap(h, bootstrapPeers)

	protocolCid, err := cidV1FromString(Protocol)
	t10sctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = dhtNode.Provide(t10sctx, protocolCid, true)
	util.CheckErr(err)

	t10sctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	peers, err := dhtNode.FindProviders(t10sctx, protocolCid)
	util.CheckErr(err)

	for _, pi := range peers {
		if err := connectPeer(h, pi); err != nil {
			log.Println(err)
			continue
		}

	}

}

func streamHandler(stream net.Stream) {
	log.Println(stream)
	stream.Reset()
}

func bootstrap(h host.Host, peers []string) {
	for _, p := range peers {
		pi, err := parsePeerInfo(p)
		if err != nil {
			log.Println(err)
			continue
		}
		if err := connectPeer(h, *pi); err != nil {
			log.Println(err)
		}
	}
}

func connectPeer(h host.Host, pi peerstore.PeerInfo) error {
	t10sctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return h.Connect(t10sctx, pi)
}

func parsePeerInfo(str string) (pi *peerstore.PeerInfo, err error) {
	addr, err := ipfsaddr.ParseString(str)
	if err != nil {
		return
	}

	pi, err = peerstore.InfoFromP2pAddr(addr.Multiaddr())
	if err != nil {
		return
	}
	return
}

func cidV1FromString(s string) (*cid.Cid, error) {
	return cid.V1Builder{Codec: cid.Raw, MhType: multihash.SHA3_256, MhLength: 0}.Sum([]byte(s))
}
