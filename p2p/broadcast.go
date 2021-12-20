package p2p

import (
	"encoding/json"
	"fmt"
	"log"
)

type BroadcastMessage struct {
	Type    string
	Content []byte
	From    string
}

func (node *Node) broadcastNewBlock(prev string, block []byte) {
	msg := &BroadcastMessage{
		Type:    "BLOCK",
		Content: block,
		From:    nodeIPAddress,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}
	node.broadcast("newblock", prev, payload)
}

func (node *Node) broadcastBlockchain(payload []byte) {
	// TODO update longest chain in the net
}

func (node *Node) broadcast(cmd string, prev string, payload []byte) {
	fmt.Printf("[BROADCAST] Broadcast triggered by %s\n", prev)
	request := append(commandToBytes(cmd), payload...)
	for _, peer := range node.peers {
		if peer != nodeIPAddress && peer != prev {
			fmt.Printf("[BROADCAST] Broadcast '%s' to %s\n", cmd, peer)
			err := sendData(peer, request)
			if err != nil {
				node.removePeer(peer)
			}
		}
	}
}
