package config

import (
	"encoding/json"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"io/ioutil"
	"os"
)

const (
	DefaultPort int32 = 22330
)

type TaggedPeerInfo struct {
	Tag   string
	ID    string
	Addrs []string
}

type Config struct {
	PrivKey string
	Port    int32
	Peers   []TaggedPeerInfo
}

func SaveConfigByData(h host.Host, privKey crypto.PrivKey) error {
	privKeyData, err := privKey.Bytes()
	if err != nil {
		return err
	}

	c := Config{
		PrivKey: crypto.ConfigEncodeKey(privKeyData),
		Port:    DefaultPort,
	}

	return SaveConfig(&c)
}

func SaveConfig(c *Config) error {
	configData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("config.json", configData, os.FileMode(0755))
	return err
}

func ReadConfig() (*Config, error) {
	fData, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var c Config
	err = json.Unmarshal(fData, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
