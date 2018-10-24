package blinkp2p

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	config "github.com/blinkspark/golab/blink-config"
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
		libp2p.EnableRelay(0, 1, 2),
		libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		return nil, nil, err
	}
	return h, privKey, nil
}

func NewHost(path string, port int32) (host.Host, error) {
	h, privKey, err := makeRandomHost(port)
	if err != nil {
		return nil, err
	}
	c := config.NewConfig()

	privKeyData, err := privKey.Bytes()
	if err != nil {
		return nil, err
	}
	c.Set("PrivateKey", crypto.ConfigEncodeKey(privKeyData))
	c.Set("Port", port)
	c.Save(path)
	return h, nil
}

func NewHostFromConfig(path string) (host.Host, error) {
	c, err := config.LoadConfig(path)
	if err != nil {
		return nil, err
	}
	privKeyStr, ok := c.Get("PrivateKey").(string)
	if !ok {
		return nil, errors.New("PrivateKey not found")
	}
	port, ok := c.Get("Port").(float64)
	if !ok {
		return nil, errors.New("Port not found")
	}

	privKeyData, err := crypto.ConfigDecodeKey(privKeyStr)
	if err != nil {
		return nil, err
	}

	privKey, err := crypto.UnmarshalPrivateKey(privKeyData)
	if err != nil {
		return nil, err
	}

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", int(port))

	return libp2p.New(context.Background(),
		libp2p.NATPortMap(),
		libp2p.Identity(privKey),
		libp2p.EnableRelay(0, 1, 2),
		libp2p.ListenAddrStrings(listenAddr),
	)
}
