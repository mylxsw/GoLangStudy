package main

import "github.com/lunny/tango"

func main() {
	t := tango.Classic()
	t.Get("/", func() string {
		return "Hello, tango!"
	})

	t.Run()
}
