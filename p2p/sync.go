package p2p

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func sendHello(addr string) {
	payload := nodeAddress
	request := append(commandToBytes("hello"), payload...)

	sendData(addr, request)
}

func sendConnect(from string, to string) {
	request := append(commandToBytes("connect"), to...)
	sendData(from, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		removeNAKnownNode(addr)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}