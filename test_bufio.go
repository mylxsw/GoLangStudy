package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, _ := os.Open("/Users/mylxsw/codes/github/mylxsw.github.io/_posts/2014-10-31-PHP扩展开发(一)构建第一个扩展.md")
	defer file.Close()

	var metas map[string]string = make(map[string]string, 10)

	scaner := bufio.NewScanner(file)
	scaner.Scan()
	if scaner.Text() == "---" {
		for scaner.Scan() {
			if scaner.Text() == "---" {
				break
			}

			items := strings.SplitN(scaner.Text(), ":", 2)
			metas[items[0]] = items[1]
		}
	}

	for key, val := range metas {
		fmt.Printf("%s  =>  %s\n", key, val)
	}

	fmt.Println("\n正文部分\n")

	content := ""
	for scaner.Scan() {
		content += scaner.Text() + "\n"
	}

	fmt.Println(content)
}
