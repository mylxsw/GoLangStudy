package main

import (
	"fmt"
	"log"
	"sync"
	"aicode.cc/remote-tail/command"
	"aicode.cc/remote-tail/console"
	"flag"
	"strings"
)

var welcomeMessage string = `
 ____                      _      _____     _ _
|  _ \ ___ _ __ ___   ___ | |_ __|_   _|_ _(_) |
| |_) / _ \ '_ ' _ \ / _ \| __/ _ \| |/ _' | | |
|  _ <  __/ | | | | | (_) | ||  __/| | (_| | | |
|_| \_\___|_| |_| |_|\___/ \__\___||_|\__,_|_|_|

author: mylxsw
homepage: github.com/mylxsw/remote-tail
` + "\x1b[0;31m-----------------------------------------------\x1b[0m\n"

var filepath *string = flag.String("file", "/var/log/messages", "-file=\"/home/data/logs/**/*.log\"")
var hostStr *string = flag.String("hosts", "root@127.0.0.1", "-hosts=root@192.168.1.225,root@192.168.1.226")

func main() {
	flag.Parse()

	var hosts []string = strings.Split(*hostStr, ",")
	var script string = fmt.Sprintf("tail -f %s", *filepath)

	fmt.Println(welcomeMessage)

	outputs := make(chan command.Message, 20)
	var wg sync.WaitGroup

	for _, host := range hosts {
		wg.Add(1)
		go func(host, script string) {
			cmd, err := command.NewCommand(host, script)
			if err != nil {
				log.Fatal(err)
			}

			cmd.Execute(outputs, nil)

			wg.Done()
		}(host, script)
	}

	wg.Add(1)
	go func() {
		for output := range outputs {
			fmt.Printf(
				"%s %s %s",
				console.ColorfulText(console.TextGreen, output.Host),
				console.ColorfulText(console.TextYellow, "->"),
				output.Content,
			)
		}

		wg.Done()
	}()

	wg.Wait()
}
