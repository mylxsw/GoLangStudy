package main

import (
	"fmt"
	"strconv"
)

func main() {
	a()
	b()
	fmt.Println(c())
}

func a() {
	i := 0
	defer fmt.Println("defer后的函数参数在调用 defer 时候就已经确定值了: " + strconv.Itoa(i))
	i++
	return
}

func b() {
	defer fmt.Println("")
	for i := 0; i < 4; i++ {
		defer fmt.Print(i)
	}
	fmt.Print("FILO:")
}

func c() (i int) {
	defer func() {
		i++
	}()

	return 1
}
