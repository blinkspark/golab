package config

import (
	"encoding/json"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-peer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"os"
)

const (
	DefaultPort int32 = 22330
)

type TaggedPeerInfo struct {
	ps.PeerInfo
	Tag string
}

func (pi *TaggedPeerInfo) MarshalJSON() ([]byte, error) {
	out := make(map[string]interface{})
	out["ID"] = pi.ID.Pretty()
	var addrs []string
	for _, a := range pi.Addrs {
		addrs = append(addrs, a.String())
	}
	out["Addrs"] = addrs
	out["Tag"] = pi.Tag
	return json.Marshal(out)
}

func (pi *TaggedPeerInfo) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	pid, err := peer.IDB58Decode(data["ID"].(string))
	if err != nil {
		return err
	}
	pi.ID = pid
	addrs, ok := data["Addrs"].([]interface{})
	if ok {
		for _, a := range addrs {
			pi.Addrs = append(pi.Addrs, multiaddr.StringCast(a.(string)))
		}
	}
	pi.Tag = data["Tag"].(string)
	return nil
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
