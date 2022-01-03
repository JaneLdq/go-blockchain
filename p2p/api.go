package p2p

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

var nodeIPAddress string
var node Node

/* When a node is started, it is isolated.
 * Use the CLI command `connect` to connect two nodes,
 * the first time two nodes conneted, they will sync lastest chain,
 * and also registered in each other's known peer for further data sync actions
 */
func StartNode(nodeId uint) {
	node = Node{}
	err := node.loadFromFile(nodeId)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Node %d not found.", nodeId))
	}
	fmt.Printf("[API] start node %d, whose blockchain height is %d\n", nodeId, node.getHeight())

	nodeIPAddress = fmt.Sprintf("localhost:%d", nodeId)
	ln, err := net.Listen(PROTOCOL, nodeIPAddress)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Node %d failed to start.", nodeId))
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConn(conn)
	}
}

func NewNode(port uint) *Node {
	nodeFile := buildNodeFilePath(port)
	if _, err := os.Stat(nodeFile); err == nil {
		return nil
	}
	private, public := generateKeyPair()
	addr := generateAddress(public)
	node := Node{
		Port:       port,
		Address:    addr,
		PrivateKey: private,
		Publickey:  public,
		peers:      []string{},
	}
	node.saveToFile()
	return &node
}

// return the base58 Address of a node
func GetAddress(nodeId uint) string {
	node := Node{}
	err := node.loadFromFile(nodeId)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Node %d not found.", nodeId))
	}
	return node.Address
}

// from and to are both IP address, e.g, localhost:3000
func ConnectNode(from string, to string) {
	request := append(CONNECT.ToByteArray(), to...)
	sendData(from, request)
}

type SendMessage struct {
	From   string
	To     string
	Amount string
}

func Mine(nodeIpAddr string, from string, to string, amount string) {
	msg := &SendMessage{
		From:   from,
		To:     to,
		Amount: amount,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	sendData(nodeIpAddr, append(MINE.ToByteArray(), payload...))
}
