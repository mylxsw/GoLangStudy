package main

import (
	"log"
	"net"

	"aicode.cc/2016/socket/protocol"
)

func reader(readerChan chan []byte) {
	for {
		select {
		case data := <-readerChan:
			log.Print(string(data))
		}
	}
}

func handleConnection(conn net.Conn) {
	tmpBuffer := make([]byte, 0)
	readerChan := make(chan []byte, 16)

	go reader(readerChan)

	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("%s connection error: %v", conn.RemoteAddr().String(), err)
			return
		}

		// log.Printf("%s receive data: %v", conn.RemoteAddr().String(), string(buffer[:n]))
		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChan)
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}()

	netListen, err := net.Listen("tcp", "localhost:10240")
	if err != nil {
		panic(err)
	}
	defer netListen.Close()

	log.Print("Waiting for client...")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		log.Printf("%s tcp connect success", conn.RemoteAddr().String())
		handleConnection(conn)
	}
}
