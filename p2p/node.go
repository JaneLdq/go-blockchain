package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

var nodeIPAddress string
var node Node

// TODO 
var blocks []string

type Node struct {
	Port       uint
	Address    string
	PrivateKey ecdsa.PrivateKey
	Publickey  []byte
	peers      []string
}

func (node *Node) isBlockExisted(newBlock string) bool {
	for _, block := range blocks {
		if block == newBlock {
			return true
		}
	}
	return false
}

func (node *Node) broadcast(from string, payload []byte) {
	request := append(commandToBytes("update"), payload...)
	for _, peer := range node.peers {
		if peer != nodeIPAddress && peer != from {
			fmt.Printf("[NODE] Broadcast to %s\n", peer)
			err := sendData(peer, request)
			if err != nil {
				node.removePeer(peer)
			}
		}
	}
}

func (node *Node) addPeer(newPeer string) {
	for _, peer := range node.peers {
		if (peer == newPeer) {
			return
		}
	}
	node.peers = append(node.peers, newPeer)
	fmt.Printf("New peer added: %s\n", newPeer)
}

func (node *Node) removePeer(lostPeer string) {
	var peerIdx int
	for idx, peer := range node.peers {
		if (peer == lostPeer) {
			peerIdx = idx
			break
		}
	}
	node.peers = append(node.peers[0:peerIdx], node.peers[peerIdx+1:]...)
	fmt.Printf("Peer removed: %s\n", lostPeer)
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

func ValidateNodeAddress(address string) bool {
	pubKeyHash := base58.Decode(address)
	actualChecksum := pubKeyHash[len(pubKeyHash)-ADDR_CHECKSUM_LEN:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-ADDR_CHECKSUM_LEN]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Equal(actualChecksum, targetChecksum)
}
