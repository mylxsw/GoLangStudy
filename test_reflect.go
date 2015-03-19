package main

import (
	"fmt"
	"log"
	"reflect"
)

func CallUserFunc(function interface{}) {
	reflectValue := reflect.ValueOf(function)
	if reflectValue.Kind() != reflect.Func {
		panic("参数必须为函数类型!")
	}
	log.Println("开始调用函数...")
	reflectValue.Call(nil)
	log.Println("函数调用完成.")
}

type TestObject struct {
	username string
	age      int
}

func (this *TestObject) SetAge(age int) {
	this.age = age
}

func (this *TestObject) interMethod() {
	fmt.Println("Intermethod")
}

func main() {
	var test_func = func() {
		fmt.Println("Hello, function")
	}

	val := reflect.ValueOf(test_func)
	fmt.Println(val.Kind())

	if val.Kind() == reflect.Func {
		fmt.Println("test_func 是函数类型")
	}

	CallUserFunc(test_func)

	reflect_object := reflect.ValueOf(&TestObject{})
	fmt.Printf("Method count: %d\n", reflect_object.Type().NumMethod())
	fmt.Printf("Name: %s\n", reflect_object.Type().Name)

	allowMethods := make(map[string]bool)
	for i := 0; i < reflect_object.Type().NumMethod(); i++ {
		method := reflect_object.Type().Method(i)
		fmt.Println(method.Name)
		if method.Name[0] >= 'A' && method.Name[0] <= 'Z' {
			allowMethods[method.Name] = true
		} else {
			allowMethods[method.Name] = false
		}
	}

	fmt.Println(allowMethods)

}
