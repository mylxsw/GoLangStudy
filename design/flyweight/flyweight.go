package flyweight

import (
	"fmt"
	"strconv"
)

type Flyweight interface {
	Operation(extrinsicstate int)
}

type ConcreteFlyweight struct{}

func (fly *ConcreteFlyweight) Operation(extrinsicstate int) {
	fmt.Println("具体Flyweight:" + strconv.Itoa(extrinsicstate))
}

type UnsharedConcreteFlyweight struct{}

func (fly *UnsharedConcreteFlyweight) Operation(extrinsicstate int) {
	fmt.Println("不同享的具体Flyweight:" + strconv.Itoa(extrinsicstate))
}

type FlyweightFactory struct {
	flyweight map[string]Flyweight
}

func GetFlyweightFactory() *FlyweightFactory {
	factory := &FlyweightFactory{}
	factory.flyweight = make(map[string]Flyweight)
	factory.flyweight["X"] = &ConcreteFlyweight{}
	factory.flyweight["Y"] = &ConcreteFlyweight{}
	factory.flyweight["Z"] = &ConcreteFlyweight{}

	return factory
}

func (factory *FlyweightFactory) GetFlyweight(key string) Flyweight {
	return factory.flyweight[key]
}
