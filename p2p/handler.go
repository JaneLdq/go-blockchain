package p2p

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
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
	case "newblock":
		node.handleNewBlock(request)
	case "reqchain":
		node.handleReqChain(request)
	case "updatechain":
		node.handleUpdateChain(request)
	default:
		log.Fatalln("Unknown command!")
	}
	conn.Close()
}

type HelloMessage struct {
	From             string
	Address          string
	Height           uint
	LastestBlockHash string
}

func (node *Node) handleConnect(request []byte) {
	payload := getPayload(request)
	newPeer := string(payload)

	fmt.Printf(logTemp, "connect", newPeer)

	msg := HelloMessage{
		From:             nodeIPAddress,
		Address:          node.Address,
		Height:           node.Height,
		LastestBlockHash: "",
	}

	data, _ := json.Marshal(msg)

	// when connect to a peer, the calling node should sync its blockchain with the other node
	err := sendData(newPeer, append(commandToBytes("hello"), data...))
	if err == nil {
		node.addPeer(newPeer)
	}
}

func (node *Node) handleHello(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, "hello", payload)

	msg := HelloMessage{}
	json.Unmarshal(payload, &msg)

	node.addPeer(msg.From)

	if msg.Height > node.Height {
		// if local blockchain is shorter, request blockchain from the new peer and broadcast to other known peers
		sendData(msg.From, append(commandToBytes("reqchain"), []byte(nodeIPAddress)...))
	} else if msg.Height < node.Height {
		// if local blockchain is longer, send local blockchain to the new peer
		node.sendChain(msg.From)
	}
}

func (node *Node) handleMine() {
	fmt.Printf(logTemp, "mine", "")

	// TODO mining
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	block := "dummyblock" + strconv.Itoa(r1.Intn(1000))
	// TODO add block to local block chain
	blocks = append(blocks, block)
	node.broadcastNewBlock(nodeIPAddress, []byte(block))
}

func (node *Node) handleNewBlock(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, "newblock", payload)

	msg := BroadcastMessage{}
	json.Unmarshal(payload, &msg)

	// TODO check if the block already exist
	existed := node.isBlockExisted(msg.Content)

	// broadcast to known peers if this is a new block
	if !existed {
		blocks = append(blocks, string(msg.Content))
		node.broadcastNewBlock(msg.From, msg.Content)
	} else {
		fmt.Printf("[HANDLER] Drop Broadcast Message: %s\n", payload)
	}
}

func (node *Node) handleReqChain(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, "reqchain", payload)
	peer := string(payload)
	node.sendChain(peer)
}

func (node *Node) handleUpdateChain(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, "updatechain", payload)
	// TODO update chain
}

func (node *Node) sendChain(destination string) {
	// TODO get chain from database
	sendData(destination, append(commandToBytes("updatechain"), "this is latest chain"...))
}