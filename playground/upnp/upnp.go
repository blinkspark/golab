package main

import (
	"fmt"

	"github.com/blinkspark/golab/util"
	"github.com/prestonTao/upnp"
)

func main() {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(62233, 62233, "TCP"); err == nil {
		fmt.Println("success !")
		// remove port mapping in gatway

		// http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		// 	res.Write([]byte("Hello"))
		// })
		// http.ListenAndServe(":62233", nil)

		mapping.Reclaim()
	} else {
		fmt.Println("fail !")
	}
}

func upnpTest() {
	mapping := new(upnp.Upnp)
	err := mapping.ExternalIPAddr()
	fmt.Println(mapping.GatewayOutsideIP)
	util.CheckErr(err)
	if err := mapping.AddPortMapping(55789, 55789, "TCP"); err == nil {
		fmt.Println("success !")
		// remove port mapping in gatway
		mapping.Reclaim()
	} else {
		fmt.Println("fail !")
	}
}
