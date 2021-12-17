package cmd

import (
	"fmt"
	blc "go-blockchain/blc"
	p2p "go-blockchain/p2p"
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

func (cli *CLI) connectNodes(from string, to string) {
	p2p.ConnectNode(from, to)
}

func (cli *CLI) startNode(nodeId uint) {
	p2p.StartNode(nodeId)
}

func (cli *CLI) initChain(nodeId uint) {
	// TODO init chain
}

func (cli *CLI) address(nodeId uint) {
	// TODO get node address
}

func (cli *CLI) send(from string, to string, message string) {
	// TODO send message from one node to another (trigger transaction)
}

func (cli *CLI) createBlockchain(address string) {

	bc := blc.CreateBlockchain(address)
	defer bc.DB.Close()

	fmt.Println("Done!")
}
