package p2p

import (
	"encoding/json"
	"fmt"
	"go-blockchain/blc"
	"io/ioutil"
	"log"
	"net"
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
		node.handleMine(request)
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
	Height        int
	LastBlockHash string
}

func (node *Node) handleConnect(request []byte) {
	payload := getPayload(request)
	newPeer := string(payload)

	fmt.Printf(logTemp, CONNECT.String(), newPeer)

	msg := HelloMessage{
		From:          nodeIPAddress,
		Address:       node.Address,
		Height:        node.getHeight(),
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

	height := node.getHeight()

	if msg.Height > height {
		// if local blockchain is shorter, request blockchain from the new peer and broadcast to other known peers
		sendData(msg.From, append(REQ_CHAIN.ToByteArray(), []byte(nodeIPAddress)...))
	} else if msg.Height < height {
		// if local blockchain is longer, send local blockchain to the new peer
		node.sendChain(msg.From)
	}
}


func (node *Node) handleMine(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, MINE.String(), payload)

	msg := SendMessage{}
	json.Unmarshal(payload, &msg)

	bc := blc.NewBlockchainWithGenesis(node.Port)
	defer bc.DB.Close()

	bc.MineNewBlock(msg.From, msg.To, msg.Amount)
	block := bc.Iterator().Next()
	node.broadcastNewBlock(nodeIPAddress, block.Serialize())
}

func (node *Node) handleNewBlock(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, NEW_BLOCK.String(), payload)

	msg := BroadcastMessage{}
	json.Unmarshal(payload, &msg)

	block := blc.DeserializeBlock(msg.Content)
	valid := node.isBlockValid(*block)

	// broadcast to known peers if this is a new block
	if valid {
		blc.AddNewBlock(node.Port, *block)
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


func (node *Node) handleUpdateChain(request []byte) {
	payload := getPayload(request)
	fmt.Printf(logTemp, UPDATE_CHAIN.String(), "")

	msg := BroadcastMessage{}
	json.Unmarshal(payload, &msg)

	chain := blc.BlockchainOjbect{}
	json.Unmarshal(msg.Content, &chain)

	if chain.Height > node.getHeight() {
		fmt.Printf("[HANDLER] Replace with longer chain (Height: %d).\n", chain.Height)
		blc.UpdateChain(node.Port, chain.Blocks)
		node.broadcastChain(msg.From, []byte(msg.Content))
	} else {
		fmt.Printf("[HANDLER] Outdated chain, ignored")
	}
}

func (node *Node) sendChain(destination string) {
	chain := blc.GetChain(node.Port)
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
