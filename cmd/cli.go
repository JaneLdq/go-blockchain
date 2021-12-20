package cmd

import (
	"fmt"
	blc "go-blockchain/blc"
	p2p "go-blockchain/p2p"
)

type CLI struct{}

func (cli *CLI) newNode(port uint) {
	node := p2p.NewNode(port)
	if node == nil {
		fmt.Printf("Node running on port %d existed, please choose another port.\n", port)
		return
	}
	fmt.Printf("New node %d created with address %s\n", port, node.Address)
	bc := blc.CreateBlockchain(node.Address)
	defer bc.DB.Close()
	fmt.Printf("New blockchain initialized on node %d!\n", port)
}

func (cli *CLI) connectNodes(from string, to string) {
	p2p.ConnectNode(from, to)
}

func (cli *CLI) startNode(nodeId uint) {
	p2p.StartNode(nodeId)
}

func (cli *CLI) address(nodeId uint) {
	fmt.Printf("Node #%d address: %s\n", nodeId, p2p.GetAddress(nodeId))
}

func (cli *CLI) send(from string, to string, amount string) {
	// TODO send coins from an address to another
}

func (cli *CLI) mine(nodeAddr string) {
	p2p.Mine(nodeAddr)
}

func (cli *CLI) createChain(address string) {
	bc := blc.CreateBlockchain(address)
	defer bc.DB.Close()
	fmt.Println("Done!")
}

func (cli *CLI) printChain(nodeId uint) {
	blc.PrintChain(nodeId)
}
