package p2p

import (
	"math/rand"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"
	"fmt"
)

const logTemp = "[HANDLER] Received '%s' with payload = {%s}\n"

func handleConn(conn net.Conn) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:P2P_CMD_LEN])
	switch command {
	case "connect":
		node.handleConnect(request)
	case "hello":
		node.handleHello(request)
	case "mine":
		node.handleMine()
	case "update":
		node.handleUpdate(request)
	default:
		log.Fatalln("Unknown command!")
	}
	conn.Close()
}

func (node *Node) handleConnect(request []byte) {
	buff := getPayload(request)
	newPeer := buff.String()

	fmt.Printf(logTemp, "connect", newPeer)

	// send a hello message to the other peer
	err := sendData(newPeer, append(commandToBytes("hello"), []byte(nodeIPAddress)...))
	if err == nil {
		node.addPeer(newPeer)
	}
}

func (node *Node) handleHello(request []byte) {
	buff := getPayload(request)
	newPeer := buff.String()

	fmt.Printf(logTemp, "hello", newPeer)

	node.addPeer(newPeer)
}

type UpdateMessage struct {
	Type string
	Content string
	From string
}

func (node *Node) handleMine() {
	fmt.Printf(logTemp, "mine", "")

	// TODO mining
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
	block := "dummyblock" + strconv.Itoa(r1.Intn(1000))
	// TODO add block to local block chain
	blocks = append(blocks, block)

	msg := &UpdateMessage{
		Type: "BLOCK",
		Content: block,
		From: nodeIPAddress,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}
	node.broadcast(nodeIPAddress, payload)
}

func (node *Node) handleUpdate(request []byte) {
	buff := getPayload(request)

	fmt.Printf(logTemp, "update", buff.Bytes())

	msg := UpdateMessage{}
	json.Unmarshal(buff.Bytes(), &msg)

	// TODO check if the block already exist
	existed := node.isBlockExisted(msg.Content)

	// broadcast to known peers if this is a new block
	if !existed {
		blocks = append(blocks, msg.Content)
		from := msg.From
		msg.From = nodeIPAddress
		relayMsg, err := json.Marshal(msg)
		if err != nil {
			log.Panic(err)
		}
		node.broadcast(from, relayMsg)
	} else {
		fmt.Printf("[HANDLER] Dropped UpdateMessage: %s\n", buff.Bytes())
	}
}