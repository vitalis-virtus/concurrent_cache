package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	v1 "github.com/vitalis-virtus/concurrent_cache/memo/v1"
	v2 "github.com/vitalis-virtus/concurrent_cache/memo/v2"
)

var incomingURLs = []string{
	"https://golang.org",
	"https://godoc.org",
	"https://play.golang.org",
	"http://gopl.io",
	"https://golang.org",
	"https://godoc.org",
	"https://play.golang.org",
	"http://gopl.io",
}

func main() {
	memCacheV1 := v1.New(httpGetBody)

	var nV1 sync.WaitGroup
	for _, url := range incomingURLs {
		nV1.Add(1)
		go func(url string) {
			start := time.Now()
			value, err := memCacheV1.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d байтов\n", url, time.Since(start), len(value.([]byte)))
			nV1.Done()
		}(url)
	}
	nV1.Wait()

	memCacheV2 := v2.New(httpGetBody)

	var nV2 sync.WaitGroup
	for _, url := range incomingURLs {
		nV2.Add(1)
		go func(url string) {
			start := time.Now()
			value, err := memCacheV2.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d байтов\n", url, time.Since(start), len(value.([]byte)))
			nV2.Done()
		}(url)
	}
	nV2.Wait()
}

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
