package main

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/blinkspark/golab/util"
	mycrypto "github.com/blinkspark/golab/util/crypto"
)

func main() {
	password := []byte("123456")
	var plainData []byte

	fname := "playground/crypto_test/test.txt"
	f, err := os.Open(fname)
	util.CheckErr(err)

	fr := bufio.NewReader(f)
	plainData, err = ioutil.ReadAll(fr)
	util.CheckErr(err)

	ctext, err := mycrypto.CFB_Enc(password, plainData)
	util.CheckErr(err)
	// fmt.Println((ctext))
	util.Ignore(ctext)

	t, err := mycrypto.CFB_Dec(password, ctext)
	// fmt.Println(string(t))
	util.Ignore(t)

	encf, err := os.Create("playground/crypto_test/test.enc")
	util.CheckErr(err)

	encfw := bufio.NewWriter(encf)
	_, err = encfw.Write(ctext)
	util.CheckErr(err)
}

func unifyLen(target []byte, min int) (result []byte) {
	tlen := len(target)
	if tlen < min {
		pad := make([]byte, min-tlen)
		target = append(target, pad...)
	} else if tlen > min {
		target = target[:min]
	}
	return target
}
