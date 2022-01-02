package p2p

import (
	"encoding/json"
	"fmt"
	"log"
)

type MessageType uint8

const (
	BLOCK MessageType = iota
	CHAIN
)
type BroadcastMessage struct {
	Type    MessageType
	Content []byte
	From    string
}

func (node *Node) broadcastNewBlock(prev string, block []byte) {
	msg := &BroadcastMessage{
		Type:    BLOCK,
		Content: block,
		From:    nodeIPAddress,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}
	node.broadcast(NEW_BLOCK, prev, payload)
}

func (node *Node) broadcastChain(prev string, payload []byte) {
	// TODO braodcast longest chain in the net
	msg := &BroadcastMessage{
		Type:    CHAIN,
		Content: payload,
		From:    nodeIPAddress,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}
	node.broadcast(UPDATE_CHAIN, prev, payload)
}

func (node *Node) broadcast(cmd CMD, prev string, payload []byte) {
	fmt.Printf("[BROADCAST] Broadcast triggered by %s\n", prev)

	request := append(cmd.ToByteArray(), payload...)
	for _, peer := range node.peers {
		if peer != nodeIPAddress && peer != prev {
			fmt.Printf("[BROADCAST] Broadcast '%s' to %s\n", cmd.String(), peer)

			err := sendData(peer, request)
			if err != nil {
				node.removePeer(peer)
			}
		}
	}
}
