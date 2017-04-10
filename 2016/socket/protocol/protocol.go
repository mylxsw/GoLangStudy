package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	Header         = "aicode.cc"
	HeaderLength   = len(Header)
	SaveDataLength = 4
)

func Packet(message []byte) []byte {
	return append(append([]byte(Header), IntToBytes(len(message))...), message...)
}

func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+HeaderLength+SaveDataLength {
			break
		}
		if string(buffer[i:i+HeaderLength]) == Header {
			messageLength := BytesToInt(buffer[i+HeaderLength : i+HeaderLength+SaveDataLength])
			if length < i+HeaderLength+SaveDataLength+messageLength {
				break
			}
			data := buffer[i+HeaderLength+SaveDataLength : i+HeaderLength+SaveDataLength+messageLength]
			readerChannel <- data

			i += HeaderLength + SaveDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)

	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
