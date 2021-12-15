package clitest

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	isValidArgs()
	//flagString := flag.String("printchain", "", "输出所有的区块信息")
	//
	//flagInt := flag.Int("number", 7, "输出一个整数")
	//
	//flagBool := flag.Bool("open", false, "判断真假。。")
	//
	//flag.Parse()
	//
	//fmt.Printf("%s\n", *flagString)
	//fmt.Printf("%d\n", *flagInt)
	//fmt.Printf("%v\n", *flagBool)

	//args := os.Args
	//
	//fmt.Printf("%v\n", args)

	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "sfqer", "交易数据")

	switch os.Args[1] {
	case "addBlock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}

		fmt.Println(*flagAddBlockData)
	}

	if printChainCmd.Parsed() {
		fmt.Println("输出所有区块的数据。。。")
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddblock -data DATA - 交易数据")
	fmt.Println("\tprintchain - 输出区块信息")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
