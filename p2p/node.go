package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

var knownNodes = []string{
	GENESIS_NODE,
}

var nodeAddress string

type Node struct {
	Port       uint
	Address    string
	Miner      bool
	PrivateKey ecdsa.PrivateKey
	Publickey  []byte
}

func StartNode(nodeId uint) {
	node := Node{}
	err := node.loadFromFile(nodeId)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Node %d not found.", nodeId))
	}

	nodeAddress = fmt.Sprintf("localhost:%d", nodeId)
	ln, err := net.Listen(PROTOCOL, nodeAddress)
	if err != nil {
		log.Fatal(fmt.Sprintf("Node %d failed to start.", nodeId))
	}
	defer ln.Close()

	// TODO init block chain from db
	// bc := NewBlockchain(nodeID)

	// register in genesis node
	if nodeAddress != GENESIS_NODE {
		sendHello(GENESIS_NODE)
		// TODO update chain
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConn(conn)
	}
}

func GetAddress(nodeId uint) string {
	node := Node{}
	err := node.loadFromFile(nodeId)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Node %d not found.", nodeId))
	}
	return node.Address
}

func ConnectNode(addrFrom string, addrTo string) {
	sendConnect(addrFrom, addrTo)
}

func NewNode(port uint, miner bool) *Node {
	nodeFile := buildNodeFilePath(port)
	if _, err := os.Stat(nodeFile); err == nil {
		return nil
	}
	private, public := generateKeyPair()
	addr := generateAddress(public)
	node := Node{port, addr, miner, private, public}
	node.saveToFile()
	return &node
}

func generateKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

/**
 * 	The process to generate address:
 *  1. pubKeyHash = ripemd160(sha256(pubKey)
 *  2. checksum = sha256(sha256(0x00 + pubKeyHash))
 *  3. address = base58(0x00 + pubKeyHash + checksum)
 */
func generateAddress(pubKey []byte) string {
	pubKeyHash := btcutil.Hash160(pubKey)
	versionedPayload := append([]byte{ADDR_TYPE}, pubKeyHash...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	return base58.Encode(fullPayload)
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:ADDR_CHECKSUM_LEN]
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}
	return false
}

func addKnownNode(addr string) {
	if !nodeIsKnown(addr) {
		knownNodes = append(knownNodes, addr)
	}
}

func removeNAKnownNode(addr string) {
	var updatedNodes []string
	for _, node := range knownNodes {
		if node != addr {
			updatedNodes = append(updatedNodes, node)
		}
	}
	knownNodes = updatedNodes
}

func ValidateNodeAddress(address string) bool {
	pubKeyHash := base58.Decode(address)
	actualChecksum := pubKeyHash[len(pubKeyHash)-ADDR_CHECKSUM_LEN:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-ADDR_CHECKSUM_LEN]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Equal(actualChecksum, targetChecksum)
}
