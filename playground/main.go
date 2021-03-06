package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	uuid "github.com/satori/go.uuid"

	"bufio"
	"log"

	"github.com/blinkspark/golab/util"
)

func main() {
	cmdTest()
}

func cmdTest() {
	cmd := exec.Command("ls")

	reader, err := cmd.StdoutPipe()
	util.CheckErr(err)
	defer reader.Close()

	err = cmd.Start()
	util.CheckErr(err)

	bufReader := bufio.NewReader(reader)
	str, err := bufReader.ReadString(0)
	if err != nil && err != io.EOF {
		log.Println(err)
		return
	}

	fmt.Println(str)

	err = cmd.Wait()
	util.CheckErr(err)
}

func uuidTest() {
	u1 := uuid.NewV4()
	fmt.Println(u1)
}

func generateKey() {
	hexKey := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, hexKey)
	util.CheckErr(err)
	str := hex.EncodeToString(hexKey)
	fmt.Println(str)
	fmt.Println(hexKey)
}

func downloadSpeed() {
	res, err := http.Get("http://i.weather.com.cn/images/cn/sjztj/2018/06/26/20180625193913993A850D318EFBC4AB12DF3F77193B5A.jpg")
	util.CheckErr(err)

	before := time.Now()
	content, err := ioutil.ReadAll(res.Body)
	util.CheckErr(err)
	bytesRead := len(content)
	after := time.Now()
	dur := after.Sub(before)
	fmt.Println("time: ", dur.Seconds())
	fmt.Println("bytes: ", bytesRead)
	BpS := float64(bytesRead) / dur.Seconds()
	KBpS := BpS / 1024
	MBpS := KBpS / 1024
	fmt.Println(MBpS, " MB/s")
}

func testRet() (err error) {
	if s, err := genRet(); err != nil {
		log.Println(s)
		log.Println(err)
	}
	return
}

func genRet() (string, error) {
	return "a", bufio.ErrBufferFull
}

type deferFunc func()

func testDefer() {
	var df deferFunc = func() {
		log.Println("a")
	}
	defer df()
	df = func() {
		log.Println("b")
	}
	defer df()
}

func testForDefer() {
	for i := 0; i < 10; i++ {
		defer func(n int) {
			log.Println(n)
		}(i)
	}
}
