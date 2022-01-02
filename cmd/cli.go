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
	bc := blc.CreateBlockchain(node.Address, port)
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

func (cli *CLI) send(from string, to string, amount string, nodeId uint) {
	bc := blc.NewBlockchainWithGenesis(nodeId)
	defer bc.DB.Close()

	bc.MineNewBlock(from, to, amount)
}

func (cli *CLI) mine(nodeIpAddr string, from string, to string, amount string) {
	p2p.Mine(nodeIpAddr, from, to, amount)
}

func (cli *CLI) printChain(nodeId uint) {
	blc.PrintChain(nodeId)
}

func (cli *CLI) getBalance(address string, nodeId uint) {
	bc := blc.NewBlockchainWithGenesis(nodeId)
	defer bc.DB.Close()
	amount := bc.GetBalance(address)
	fmt.Printf("%s have %d token\n", address, amount)
}
