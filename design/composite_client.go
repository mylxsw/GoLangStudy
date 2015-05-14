package main

import (
	comp "aicode.cc/design/composite"
	"fmt"
	"log"
)

func main() {
	fmt.Println("组合模式")

	composite := comp.NewComposite()
	composite.Draw()

	leaf := &comp.GraphLeaf{}
	leaf.Draw()

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic :", err)
		}

		composite.Add(leaf)
		comp := composite.Iterator()
		if comp != nil {
			comp.Value.Draw()
		} else {
			log.Println("迭代器为空")
		}

	}()

	leaf.Add(composite)

}
