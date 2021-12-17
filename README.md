# go-blockchain
A basic blockchain implementation in Golang.

## Installation and Usage

### CLI Usage
We provide a command line tool to setup a blockchain on your machine.

The CLI is powered by [cobra](https://github.com/spf13/cobra), which is a library for creating powerful modern CLI applications as well as a program to generate applications and command files.

run `go run . -h` to get information about the cli.

Here is a sample output, it is not a stable version, we're still woking on it.
```
Gobc is a simplified blockchain implemented in Go

Usage:
  gobc [flags]
  gobc [command]

Available Commands:
  address     Get the address of the given node
  completion  generate the autocompletion script for the specified shell
  connect     Connect a node to another
  createChain Create a blockchain and send genesis block reward to address
  help        Help about any command
  newnode     Create a new node with given port
  printChain  Print all the blocks of the blockchain
  send        Send a message from a node to another
  startnode   Start the node running on given port

Flags:
      --config string   config file (default is $HOME/.gobc.yaml)
  -h, --help            help for gobc

Use "gobc [command] --help" for more information about a command.
```
