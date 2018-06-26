package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/PuerkitoBio/goquery"

	"github.com/blinkspark/golab/http_crawler/defines"
	"github.com/blinkspark/golab/util"
	"github.com/blinkspark/golab/util/config"
)

var configEntry defines.Config
var wg sync.WaitGroup
var httpClient http.Client
var wr int32
var taskLimitChan chan int32

func main() {
	configPath := flag.String("config", "./config.json", "-config=./config.json")
	flag.Parse()
	config.InitConfig(*configPath, &configEntry)
	httpClient = makeClient(configEntry.Proxy)
	taskLimitChan = make(chan int32, configEntry.MaxRoutine)
	for _, t := range configEntry.Targets {
		beginTask(t)
	}
	wg.Wait()
}

func makeClient(urlString string) http.Client {
	client := http.Client{}
	url, err := url.Parse(urlString)
	if urlString != "" && err == nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(url),
		}
	}
	return client
}

func beginTask(t defines.Target) {
	taskLimitChan <- 0
	wg.Add(1)

	go func() {
		defer func() {
			<-taskLimitChan
			wg.Done()
		}()
		res, err := httpClient.Get(t.URL)
		if err != nil {
			fmt.Println(err)
			return
		}

		content, err := ioutil.ReadAll(res.Body)
		util.CheckErr(err)
		contentReader := bytes.NewReader(content)
		if matchSaveRule(res, t.SaveRule) {
			dirName, fileName := path.Split(res.Request.URL.Path)
			dirName = unifyDirString(configEntry.SavePath + "/" + dirName)
			err = mkDirAll(dirName)
			util.CheckErr(err)
			if fileName != "" {
				f, err := os.Create(dirName + "/" + fileName)
				util.CheckErr(err)
				defer f.Close()
				// n, err := f.Write(buffer)
				writer := bufio.NewWriter(f)

				n, err := writer.Write(content)
				util.CheckErr(err)

				err = writer.Flush()
				util.CheckErr(err)

				fmt.Println(n)
			}
		}

		doc, err := goquery.NewDocumentFromReader(contentReader)
		util.CheckErr(err)
		if t.DigRule.Deep > 0 {
			doc.Find("a").Each(func(i int, sel *goquery.Selection) {
				href, exits := sel.Attr("href")
				if exits {
					tmpURL, err := url.Parse(href)
					util.CheckErr(err)
					if tmpURL.Host == "" {
						tmpURL.Scheme = res.Request.URL.Scheme
						tmpURL.Host = res.Request.URL.Host
					}
					nextT := t
					nextT.URL = tmpURL.String()
					nextT.DigRule.Deep--
					fmt.Println(nextT)
					for atomic.LoadInt32(&wr) >= configEntry.MaxRoutine {
						// do nothing
					}

					go beginTask(nextT)
				}
			})
		}
	}()
}

func taskTest(t defines.Target) {
	wg.Add(1)
	atomic.AddInt32(&wr, 1)
	fmt.Println(atomic.LoadInt32(&wr))
	go func() {
		defer wg.Done()
		defer atomic.AddInt32(&wr, -1)
		defer fmt.Println(atomic.LoadInt32(&wr))
		res, err := httpClient.Get(t.URL)
		if err != nil {
			wg.Done()
			atomic.AddInt32(&wr, -1)
			return
		}
		// TODO retry
		// retry := configEntry.Retry
		defer res.Body.Close()
		if res.StatusCode == 200 {
			// doc, err := goquery.NewDocumentFromResponse(res)
			// util.CheckErr(err)
			content, err := ioutil.ReadAll(res.Body)
			if err != nil {
				content = make([]byte, 0)
			}
			reader := bytes.NewReader(content)
			doc, err := goquery.NewDocumentFromReader(reader)
			util.CheckErr(err)

			if matchSaveRule(res, t.SaveRule) {
				dirName, fileName := path.Split(res.Request.URL.Path)
				dirName = unifyDirString(configEntry.SavePath + "/" + dirName)
				err = mkDirAll(dirName)
				util.CheckErr(err)
				if fileName != "" {
					f, err := os.Create(dirName + "/" + fileName)
					util.CheckErr(err)
					defer f.Close()
					// n, err := f.Write(buffer)
					writer := bufio.NewWriter(f)

					n, err := writer.Write(content)
					util.CheckErr(err)

					err = writer.Flush()
					util.CheckErr(err)

					fmt.Println(n)
				}
			}

			if t.DigRule.Deep > 0 {
				doc.Find("a").Each(func(i int, sel *goquery.Selection) {
					fmt.Println("index", i)
					href, exits := sel.Attr("href")
					if exits {
						tmpURL, err := url.Parse(href)
						util.CheckErr(err)
						if tmpURL.Host == "" {
							tmpURL.Scheme = res.Request.URL.Scheme
							tmpURL.Host = res.Request.URL.Host
						}
						nextT := t
						nextT.URL = tmpURL.String()
						nextT.DigRule.Deep--
						fmt.Println(nextT)
						for atomic.LoadInt32(&wr) >= configEntry.MaxRoutine {
							// do nothing
						}

						go beginTask(nextT)
					}
				})
			}
			fmt.Println("end dig")
		}
	}()
}

func matchSaveRule(res *http.Response, saveRule defines.SaveRule) bool {
	return res.Header["Content-Type"][0] == saveRule.ContentType
}

func getURL(url string) (*http.Response, error) {
	res, err := httpClient.Get(url)
	retry := configEntry.Retry
	if err != nil && retry > 0 {
		res, err = getURL(url)
	}
	return res, err
}

func mkDirAll(dir string) error {
	var err error
	if !util.CheckDirIsExist(dir) {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	return err
}

func unifyDirString(dir string) string {
	return strings.TrimRight(strings.Replace(dir, "\\", "/", -1), " /")
}

func todo(i interface{}) {
}
