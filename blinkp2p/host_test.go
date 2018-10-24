package blinkp2p

import "testing"

func TestNewHost(t *testing.T) {
	h, err := NewHost("config.json", 22330)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(h.ID().Pretty())
		h.Close()
	}
}

func TestNewHostFromConfig(t *testing.T) {
	h, err := NewHostFromConfig("config.json")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(h.ID().Pretty())
		h.Close()
	}
}
