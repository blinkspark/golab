package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/blinkspark/golab/util"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
)

const (
	keyLength int = 4096
	testTopic     = "Test"
)

var (
	privKeyPath string
	configPath  string
	port        int
	connect     string
)

var (
	done chan bool
)

func main() {
	done = make(chan bool)
	flag.StringVar(&privKeyPath, "key", "priv.key", "-key /path/to/privatekey")
	flag.StringVar(&configPath, "c", "config.json", "-c /path/to/config.json")
	flag.IntVar(&port, "p", 22333, "-p PORT")
	flag.StringVar(&connect, "connect", "", "connect to a server: -connect /ip4/127.0.0.1/tcp/22333/ipfs/Qm***")
	flag.Parse()

	priv, err := getPrivKey(privKeyPath)
	util.CheckErr(err)

	// config, err := ReadConfigWithLazyCreate(configPath)
	// util.CheckErr(err)

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)

	host, err := libp2p.New(context.Background(),
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(listenAddr))

	ps, err := pubsub.NewFloodSub(context.Background(), host)
	util.CheckErr(err)

	if connect != "" {
		addr, err := ma.NewMultiaddr(connect)
		util.CheckErr(err)

		pi, err := peer.AddrInfoFromP2pAddr(addr)
		util.CheckErr(err)

		err = host.Connect(context.Background(), *pi)
		util.CheckErr(err)
	}

	sub, err := ps.Subscribe(testTopic)
	util.CheckErr(err)

	go handleSubscription(sub)

	if connect == "" {
		go func() {
			for {
				ps.Publish(testTopic, []byte("Hello"))
				time.Sleep(time.Second)
			}
		}()
	}

	fmt.Println(host.ID(), host.Addrs())
	select {
	case <-done:
		fmt.Println("Done")
	}
}

func handleSubscription(sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(context.Background())
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Println(msg.GetFrom(), string(msg.GetData()))
	}
}

func getPrivKey(privKeyPath string) (crypto.PrivKey, error) {
	keyData, err := ioutil.ReadFile(privKeyPath)

	// generate key if private key is not exist
	if err != nil {
		priv, _, err := crypto.GenerateKeyPair(crypto.RSA, keyLength)
		if err != nil {
			return nil, err
		}

		keyData, err = crypto.MarshalPrivateKey(priv)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(privKeyPath, keyData, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	priv, err := crypto.UnmarshalPrivateKey(keyData)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
