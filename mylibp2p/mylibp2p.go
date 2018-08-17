package mylibp2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/blinkspark/golab/mylibp2p/config"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
)

func makeRandomHost(port int32) (host.Host, crypto.PrivKey, error) {
	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)
	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	h, err := libp2p.New(context.Background(),
		libp2p.Identity(privKey),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(0, 1),
		libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		return nil, nil, err
	}
	return h, privKey, nil
}

func makeHostFromConfig(c *config.Config) (host.Host, error) {
	privKeyData, err := crypto.ConfigDecodeKey(c.PrivKey)
	if err != nil {
		return nil, err
	}
	privKey, err := crypto.UnmarshalPrivateKey(privKeyData)
	if err != nil {
		return nil, err
	}
	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", c.Port)

	return libp2p.New(context.Background(),
		libp2p.NATPortMap(),
		libp2p.Identity(privKey),
		libp2p.EnableRelay(0, 1),
		libp2p.ListenAddrStrings(listenAddr))
}

func InitHost() (host.Host, error) {
	hConfig, err := config.ReadConfig()
	if err != nil {
		newHost, privKey, err := makeRandomHost(config.DefaultPort)
		if err != nil {
			return nil, err
		}

		config.SaveConfig(newHost, privKey)
		return newHost, nil
	}
	return makeHostFromConfig(hConfig)
}
