package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
	// "error"
)

var usage = `Usage: goab [options...] <url>

Your CPU numbers: %d

`

var c = flag.Int("c", 50, "并发数目")
var n = flag.Int("n", 200, "请求数目")

func usageAndExit(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}

	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, runtime.NumCPU())
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() < 1 {
		usageAndExit("")
	}

	num := *n
	conc := *c
	url := flag.Args()[0]
	fmt.Printf("请求地址: %s\n", url)

	if num <= 0 || conc <= 0 {
		usageAndExit("n 和 c 不能小于1")
	}

	results := make(chan *result, num)
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(num)

	jobs := make(chan *http.Request, num)
	for i := 0; i < conc; i++ {
		go func() {
			worker(results, &wg, jobs)
		}()
	}

	for i := 0; i < num; i++ {
		jobs <- (&ReqOpts{
			Method: "GET",
			Url:    url,
		}).Request()
	}

	close(jobs)
	wg.Wait()

	fmt.Println("请求总时间: ", time.Now().Sub(start))
}

type result struct {
	err           error
	statusCode    int
	duration      time.Duration
	contentLength int64
}

type ReqOpts struct {
	Method string
	Url    string
}

func (r *ReqOpts) Request() *http.Request {
	req, _ := http.NewRequest(r.Method, r.Url, nil)
	return req
}

func worker(results chan *result, wg *sync.WaitGroup, ch chan *http.Request) {
	client := &http.Client{}

	for req := range ch {
		s := time.Now()
		code := 0
		size := int64(0)
		resp, err := client.Do(req)

		if err == nil {
			size = resp.ContentLength
			code = resp.StatusCode
			resp.Body.Close()
		}

		wg.Done()

		results <- &result{
			statusCode:    code,
			err:           err,
			duration:      time.Now().Sub(s),
			contentLength: size,
		}

	}
}
