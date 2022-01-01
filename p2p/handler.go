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
	command := CMD(request[0])
	switch command {
	case CONNECT:
		node.handleConnect(request)
	case HELLO:
		node.handleHello(request)
	case MINE:
		node.handleMine()
	case NEW_BLOCK:
		node.handleNewBlock(request)
	case REQ_CHAIN:
		node.handleReqChain(request)
	case UPDATE_CHAIN:
		node.handleUpdateChain(request)
	default:
		log.Fatalln("Unknown command!")
	}
	conn.Close()
}

type HelloMessage struct {
	From          string
	Address       string
	Height        uint
	LastBlockHash string
}

func (node *Node) handleConnect(request []byte) {
	payload := getPayload(request)
	newPeer := string(payload)

	fmt.Printf(logTemp, CONNECT.String(), newPeer)

	msg := HelloMessage{
		From:          nodeIPAddress,
		Address:       node.Address,
		Height:        node.Height,
		LastBlockHash: "",
	}

	data, _ := json.Marshal(msg)

	// when connect to a peer, the calling node should sync its blockchain with the other node
	err := sendData(newPeer, append(HELLO.ToByteArray(), data...))
	if err == nil {
		node.addPeer(newPeer)
	}
}

func (node *Node) handleHello(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, HELLO.String(), payload)

	msg := HelloMessage{}
	json.Unmarshal(payload, &msg)

	node.addPeer(msg.From)

	if msg.Height > node.Height {
		// if local blockchain is shorter, request blockchain from the new peer and broadcast to other known peers
		sendData(msg.From, append(REQ_CHAIN.ToByteArray(), []byte(nodeIPAddress)...))
	} else if msg.Height < node.Height {
		// if local blockchain is longer, send local blockchain to the new peer
		node.sendChain(msg.From)
	}
}

func (node *Node) handleMine() {
	fmt.Printf(logTemp, MINE.String(), "")

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
	fmt.Printf(logTemp, NEW_BLOCK.String(), payload)

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
	fmt.Printf(logTemp, REQ_CHAIN.String(), payload)
	peer := string(payload)
	node.sendChain(peer)
}

type BlockchainPayload struct {
	Height        uint
	Blocks        []byte
	LastBlockHash string
}

func (node *Node) handleUpdateChain(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, UPDATE_CHAIN.String(), payload)

	msg := BroadcastMessage{}
	json.Unmarshal(payload, &msg)

	chain := BlockchainPayload{}
	json.Unmarshal(msg.Content, &chain)

	if chain.Height > node.Height {
		node.Height = chain.Height
		fmt.Printf("[HANDLER] Updated chain height: %d\n", node.Height)
		// TODO replace local chain with received chain
		node.broadcastChain(msg.From, []byte(msg.Content))
	} else {
		fmt.Printf("[HANDLER] Outdated chain, ignored")
	}
}

func (node *Node) sendChain(destination string) {
	// TODO build chain from somewhere
	chain := &BlockchainPayload{
		Height: node.Height,
		Blocks: []byte{},
	}
	content, err := json.Marshal(chain)
	if err != nil {
		log.Panic(err)
	}

	msg := &BroadcastMessage{
		Type:    CHAIN,
		Content: content,
		From:    nodeIPAddress,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}
	sendData(destination, append(UPDATE_CHAIN.ToByteArray(), payload...))
}
