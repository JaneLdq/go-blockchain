package cmd

import (
	"fmt"
	p2p "gobc/p2p"
)

type CLI struct{}

func (cli *CLI) newNode(port uint, miner bool) {
	node := p2p.NewNode(port, miner)
	if node == nil {
		fmt.Printf("Node running on port %d existed, please choose another port.\n", port)
		return
	}
	fmt.Printf("New node %d created with address %s\n", port, node.Address)
}

func (cli *CLI) startNode(nodeId uint) {
	err := p2p.StartNode(nodeId)
	if err != nil {
		fmt.Printf("Error occurs when starting node %d\n", nodeId)
	}
}

func (cli *CLI) initChain(nodeId uint) {
	// TODO
}

func (cli *CLI) address(nodeId uint) {
	// TODO
}

func (cli *CLI) send(from string, to string, message string) {
	// TODO
}
