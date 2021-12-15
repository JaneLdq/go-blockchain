package blc

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const timeFormat string = "2006-01-02 03:04:05 PM"

type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) createBlockchain(address string) {

	bc := CreateBlockchain(address)
	defer bc.DB.Close()

	fmt.Println("Done!")
}

func (cli *CLI) printChain() {
	bc := NewBlockchainWithGenesis()
	defer bc.DB.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Block: %x\n", block.Header.Hash)
		fmt.Printf("Height: %d\n", block.Header.Height)
		fmt.Printf("Nonce: %d\n", block.Header.Nonce)
		fmt.Printf("PrevBlock: %x\n", block.Header.PrevBlockHash)
		pow := NewPoW(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println("Txs: ")
		for _, tx := range block.Txs {
			// fmt.Println(tx)
			fmt.Printf("%x\n", tx.TxHash)
			fmt.Printf("Vins: \n")
			for _, in := range tx.Vins {
				// fmt.Println(in)
				fmt.Printf("in.Txid: %x\n", in.Txid)
				fmt.Printf("n.Vout: %d\n", in.Vout)
				fmt.Printf("in.PubKey: %s\n", in.PubKey)
			}
			fmt.Printf("Vouts: \n")
			for _, out := range tx.Vouts {
				// fmt.Println(out)
				fmt.Println(out.Value)
				fmt.Println(out.PubKeyHash)
			}
		}
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Header.Timestamp, 0).Format(timeFormat))
		fmt.Printf("---------------------------------------------------------------------")
		fmt.Printf("\n\n")

		if len(block.Header.PrevBlockHash) == 0 {
			break
		}
	}
}
