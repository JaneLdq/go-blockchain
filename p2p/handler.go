package p2p

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func handleConn(conn net.Conn) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:P2P_CMD_LEN])
	
	fmt.Printf("Received '%s' command\n", command)

	switch command {
	case "hello":
		handleHello(request)
	case "connect":
		handleConnect(request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleHello(request []byte) {
	var buff bytes.Buffer
	buff.Write(request[P2P_CMD_LEN:])
	payload := buff.String()
	addKnownNode(payload)
	// for _, node := range knownNodes {
	// 	fmt.Println(node)
	// }
}

func handleConnect(request []byte) {
	var buff bytes.Buffer
	buff.Write(request[P2P_CMD_LEN:])
	payload := buff.String()
	addKnownNode(payload)
	sendHello(payload)
}