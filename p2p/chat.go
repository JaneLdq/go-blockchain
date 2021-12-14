package p2p

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

const PROTOCOL = "tcp"

var knownNodes = []string{
	"localhost:3000",
}

func StartNode(nodeId uint) error {
	node := Node{}
	err := node.loadFromFile(nodeId)
	if err != nil {
		log.Panic(err)
	}
	// start node
	nodeAddr := fmt.Sprintf("localhost:%d", nodeId)
	ln, err := net.Listen(PROTOCOL, nodeAddr)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	// bc := NewBlockchain(nodeID)

	if nodeAddr != knownNodes[0] {
		// TODO
		sendHello(knownNodes[0])
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:CTLMSG_LEN])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "hello":
		handleHello(request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleHello(request []byte) {
	var buff bytes.Buffer
	buff.Write(request[CTLMSG_LEN:])
	payload := buff.String()
	fmt.Printf("Received payload: %s\n", payload)
}

func sendHello(address string) {
	payload := "hi,there!"
	request := append(commandToBytes("hello"), payload...)

	sendData(address, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}
