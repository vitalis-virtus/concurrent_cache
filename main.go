package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	v1 "github.com/vitalis-virtus/concurrent_cache/memo/v1"
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

	var n sync.WaitGroup
	for _, url := range incomingURLs {
		n.Add(1)
		go func(url string) {
			start := time.Now()
			value, err := memCacheV1.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d байтов\n", url, time.Since(start), len(value.([]byte)))
			n.Done()
		}(url)
	}
	n.Wait()
}

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
