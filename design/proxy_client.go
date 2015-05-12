package main

import (
	"aicode.cc/design/proxy"
	"fmt"
)

func main() {
	fmt.Println("代理模式")

	proxy := proxy.NewProxy(&proxy.RealSubject{})
	proxy.Request()

}
