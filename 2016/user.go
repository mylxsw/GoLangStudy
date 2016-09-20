package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error" + err.Error())
		os.Exit(0)
	}
	fmt.Println(currentUser.Name)
	fmt.Println(currentUser.Username)
	fmt.Println(currentUser.HomeDir)
	fmt.Println(currentUser.Gid)
	fmt.Println(currentUser.Uid)

	fmt.Println(path.Join(currentUser.HomeDir, "codes/golang"))
}
