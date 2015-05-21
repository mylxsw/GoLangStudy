package main

import (
	flyweight "aicode.cc/design/flyweight"
	"fmt"
)

func main() {
	fmt.Println("享元模式")

	extrinsicstate := 32

	factory := flyweight.GetFlyweightFactory()

	fx := factory.GetFlyweight("X")
	fx.Operation(extrinsicstate)

	extrinsicstate = extrinsicstate - 1

	fy := factory.GetFlyweight("Y")
	fy.Operation(extrinsicstate)

	extrinsicstate = extrinsicstate - 1

	fz := factory.GetFlyweight("Z")
	fz.Operation(extrinsicstate)

	extrinsicstate = extrinsicstate - 1

	uf := &flyweight.UnsharedConcreteFlyweight{}
	uf.Operation(extrinsicstate)
}
