package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	//"strconv"
	"regexp"
)

type UserInfo struct {
	Username string
	Location string
	Gender   string
	Job      string
	Agree    int
	Thank    int
}

func main() {
	fmt.Println("Hello, Go")

	resp, err := http.Get("http://www.zhihu.com/people/mylxsw/followees")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	html := string(body)
}
