package adapter

import "fmt"

type Adaptee struct{}

func (ada *Adaptee) SpecificRequest() {
	fmt.Println("Specific Request")
}
