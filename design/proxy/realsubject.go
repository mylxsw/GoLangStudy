package proxy

import "fmt"

type RealSubject struct {
}

func (subject *RealSubject) Request() {
	fmt.Println("Hello, Real Subject")
}
