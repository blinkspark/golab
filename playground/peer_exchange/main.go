package main

import (
	"fmt"
	"github.com/blinkspark/golab/mylibp2p"
)

func main() {
	h, err := mylibp2p.InitHost()
	if err == nil {
		fmt.Println(h.ID().Pretty())
	}
}
