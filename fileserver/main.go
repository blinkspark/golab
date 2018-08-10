package main

import (
	"flag"
	"net/http"

	"github.com/blinkspark/golab/util"
)

func main() {
	path := flag.String("path", "", "absolute path you want to serve")
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir(*path)))
	err := http.ListenAndServe(":8080", nil)
	util.CheckErr(err)
}
