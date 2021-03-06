package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
    "sync"
    "time"
)

func main() {
	sigRecv1 := make(chan os.Signal, 1)
	sigs1 := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	fmt.Printf("Set notification for %s ... [sigRecv1]\n", sigs1)
	signal.Notify(sigRecv1, sigs1...)

	sigRecv2 := make(chan os.Signal, 1)
	sigs2 := []os.Signal{syscall.SIGQUIT}
	fmt.Printf("Set notification for %s ... [sigRecv2]\n", sigs2)

	signal.Notify(sigRecv2, sigs2...)

    var wg sync.WaitGroup
    wg.Add(2)

    go func () {
        for sig := range sigRecv1 {
            fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
        }
        fmt.Printf("End. [sigRecv1]\n")
        wg.Done()
    }()

    go func () {
        for sig := range sigRecv2 {
            fmt.Printf("Received a signal from sigRecv2: %s\n", sig)
        }

        fmt.Println("End. [sigRecv2]")
        wg.Done()
    } ()

    fmt.Println("Wait fro 2 seconds... ")
    time.Sleep(2 * time.Second)
    fmt.Println("Stop notification... ")
    signal.Stop(sigRecv1)

    close(sigRecv1)

    fmt.Println("done. [sigRecv1]")

    wg.Wait()

}
