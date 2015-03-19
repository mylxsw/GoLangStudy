package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.baidu.com",
		"http://www.letv.com",
		"http://www.csdn.net",
	}

	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			http.Get(url)
			fmt.Println("请求url: ", url)
		}(url)
	}

	wg.Wait()
}
