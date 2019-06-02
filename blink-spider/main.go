package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Job struct {
	url  string
	deep int
}

var (
	help        bool
	target      string
	sameHost    bool
	delay       int64
	depth       int
	concurrency int
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

	flag.Int64Var(&delay, "delay", 1, "delay(ms) of workers, in order to avoid blocking issues")

	flag.IntVar(&depth, "d", 2, "how deep the digger dig")
	flag.IntVar(&depth, "depth", 2, "how deep the digger dig")

	flag.IntVar(&concurrency, "c", 0, "how many goroutines")
	flag.IntVar(&concurrency, "concurrency", 0, "how many goroutines")

	flag.Parse()

	urlMap = make(map[string]bool)
	pool = make(chan Job, runtime.NumCPU()) // runtime.NumCPU()

	switch {
	case help == true:
		fmt.Println(helpText)
	case target != "null" && strings.HasPrefix(target, "http"):
		job := Job{url: target, deep: 3}
		urlMap[target] = true
		pool <- job
		tn := getThreadNum(concurrency)
		for i := 0; i < tn; i++ {
			go worker(i, pool)
		}

		timer := time.NewTimer(2 * time.Second)
		go func() {
			for {
				if len(pool) > 0 {
					timer.Reset(2 * time.Second)
				}
				time.Sleep(time.Millisecond * 100)
			}
		}()
		<-timer.C
		fmt.Println("Done")
	default:
		fmt.Println(helpText)
	}
}

func getThreadNum(c int) int {
	if c == 0 {
		return runtime.NumCPU()
	}
	return c
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

		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}
}

func handleRes(res *http.Response, doc *goquery.Document) {
	// if "text/plain" == res.Header["Content-Type"][0] {
	// dir, fileName := path.Split(res.Request.URL.Path)
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// targetDir := path.Join(currentDir, "/tmp", dir)
	// err = os.MkdirAll(targetDir, os.ModeDir)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// err = ioutil.WriteFile(path.Join(targetDir, fileName), []byte(doc.Text()), os.ModePerm)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// } // if
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
