package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/blinkspark/golab/util"
)

func main() {
	path := flag.String("path", "", "absolute path you want to serve")
	port := flag.Int("port", 8080, "port of file server")
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir(*path)))
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, nil)
	util.CheckErr(err)
}
