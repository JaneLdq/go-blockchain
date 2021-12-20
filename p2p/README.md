# P2P

main files:
* `api.go` - define public functions called by CLI
* `node.go` - define struct `Node` and functions to generate key pairs and node address (base58 address for trading transaction)
* `handler.go` - define handlers for http requests
* `broadcast.go` - define functions for broadcasting in connected nodes
* `io.go` - read & write node config from local file

## How nodes communicate with each other?
### connect
1. When a node receives a `connect` message with a peer's address, it will send a `hello` message to the peer with a payload, which contains its:
* IP
* Address
* Blockchain Height
* Latest Block Hash

2. When a node receives a `hello` message, it adds the peer as a known node and then compares the blockchain height from its own:
* If the node's local blockchain is shorter than the peer's, which means outdated, it should request the longer blockchain from the peer
2. If the node's local blockchain is longer thant the peer's, it should send its local blockchain to the peer

### mine
1. When a node successfully mined a new block, it will broadcast a `newblock` message to all its known peers.
2. When a node receives a `newblock` message, it will check whether this block is already added:
* If NOT EXIST, it adds this new block into its blockchain, and broadcast to all its known peers except the origin node
* If EXIST, it drops the block

To be continued...


