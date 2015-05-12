package proxy

import "fmt"

type Proxy struct {
	realSubject RealSubject
}

func NewProxy(subject *RealSubject) *Proxy {
	proxy := &Proxy{}
	proxy.realSubject = *subject
	return proxy
}

func (proxy *Proxy) beforeRequest() {
	fmt.Println("Before Request")
}

func (proxy *Proxy) afterRequest() {
	fmt.Println("After Request")
}

func (proxy *Proxy) Request() {
	proxy.beforeRequest()
	proxy.realSubject.Request()
	proxy.afterRequest()
}
