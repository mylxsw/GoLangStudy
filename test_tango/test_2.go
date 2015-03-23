package main

import "github.com/lunny/tango"

type Action struct {
	tango.Json
}

func (Action) Get() map[string]string {

	return map[string]string{
		"say": "Hello tango!",
	}
}

func main() {
	t := tango.Classic()
	t.Get("/", new(Action))
	t.Run()
}
