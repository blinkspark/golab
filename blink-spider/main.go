package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
	"ioutil"

	"github.com/PuerkitoBio/goquery"
)

type Job struct {
	url  string
	deep int
}

var (
	help     bool
	target   string
	sameHost bool
)

var (
	pool   chan Job
	urlMap map[string]bool
	mutex  = &sync.Mutex{}
)

const (
	helpText = `This tool is a web crawler.
Usage:
	spider -t https://google.com`
)

func main() {
	flag.BoolVar(&help, "h", false, "get help")
	flag.BoolVar(&help, "help", false, "get help")

	flag.StringVar(&target, "t", "null", "target url")
	flag.StringVar(&target, "target", "null", "target url")

	flag.BoolVar(&sameHost, "s", false, "only fetch same host")
	flag.BoolVar(&sameHost, "same", false, "only fetch same host")

	flag.Parse()

	urlMap = make(map[string]bool)
	pool = make(chan Job, runtime.NumCPU()) // runtime.NumCPU()

	switch {
	case help == true:
		fmt.Println(helpText)
	case target != "null" && strings.HasPrefix(target, "http"):
		job := Job{url: target, deep: 3}
		mutex.Lock()
		urlMap[target] = true
		mutex.Unlock()
		pool <- job
		for i := 0; i < runtime.NumCPU(); i++ {
			go worker(i, pool)
		}
		for {
			fmt.Printf("jobs len: %d\n", len(pool))
			time.Sleep(time.Second)
		}
		// select {}
	default:
		fmt.Println(helpText)
	}

}

func worker(id int, jobs chan Job) {
	for job := range jobs {
		targetURL, err := url.Parse(job.url)

		res, err := http.Get(job.url)
		if err != nil {
			log.Println(err)
			continue
		}

		doc, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			log.Println(err)
			continue
		}

		// handle response
		handleRes(res, doc)
		// dig proccess
		go dig(id, doc, targetURL, job, jobs)
	}
}

func handleRes(res *http.Response, doc *goquery.Document) {
	if "text/plain" == res.Header["Content-Type"][0] {
		dir, fileName := path.Split(res.Request.URL.Path)
		currentDir, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return
		}
		targetDir := path.Join(currentDir, dir)
		err = os.MkdirAll(targetDir, os.ModeDir)
		if err != nil {
			log.Println(err)
			return
		}

		file, err := os.Create(path.Join(targetDir, fileName))
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
	}
}

func dig(id int, doc *goquery.Document, targetURL *url.URL, job Job, jobs chan Job) {
	doc.Find("a").Each(func(i int, ele *goquery.Selection) {
		if href, exist := ele.Attr("href"); exist {
			hrefURL, err := url.Parse(href)
			if err != nil {
				log.Println(err)
				return
			}
			if hrefURL.Host == "" {
				hrefURL.Host = targetURL.Host
			}
			if hrefURL.Scheme == "" {
				hrefURL.Scheme = targetURL.Scheme
			}

			mutex.Lock()
			_, prs := urlMap[hrefURL.String()]
			mutex.Unlock()
			newDeep := job.deep - 1

			// filters
			if prs == true || newDeep < 0 {
				return
			}
			if sameHost && hrefURL.Host != targetURL.Host {
				return
			}
			if hrefURL.Scheme != "http" && hrefURL.Scheme != "https" {
				return
			}

			fmt.Println(id)
			fmt.Println(newDeep)
			fmt.Println(hrefURL)
			fmt.Println("--------")

			newJob := Job{url: hrefURL.String(), deep: newDeep}
			mutex.Lock()
			urlMap[newJob.url] = true
			mutex.Unlock()
			jobs <- newJob
		}
	})
}
