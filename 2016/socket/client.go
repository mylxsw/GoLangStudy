package main

import (
	"log"
	"net"

	"fmt"

	"aicode.cc/2016/socket/protocol"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}()

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10240")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	log.Print("connect success")

	for i := 0; i < 100; i++ {
		conn.Write(protocol.Packet([]byte(fmt.Sprintf("Hello world %d !", i))))
	}

	log.Print("send over")
}
