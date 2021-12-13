package node

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/btcsuite/btcutil/base58"
)

const nodeFile = "node_%s.dat"

type Node struct {
	Port       uint
	Address    string
	Miner      bool
	PrivateKey ecdsa.PrivateKey
	Publickey  []byte
}

func NewNode(port uint, miner bool) *Node {
	private, public := generateKeyPair()
	addr := generateAddress(public)
	node := Node{port, addr, miner, private, public}
	return &node
}

func (node *Node) LoadFromFile(nodeId string) error {
	nodeFile := fmt.Sprintf(nodeFile, nodeId)
	if _, err := os.Stat(nodeFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(nodeFile)
	if err != nil {
		log.Panic(err)
	}

	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&node)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func (node *Node) SaveToFile(nodeId string) {
	var content bytes.Buffer
	nodeFile := fmt.Sprintf(nodeFile, nodeId)

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(node)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(nodeFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func generateAddress(pubKey []byte) string {
	hash := sha256.New()
	hash.Write(pubKey)
	pubKeyHash := hash.Sum(nil)
	return base58.Encode(pubKeyHash)
}

func generateKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.Y.Bytes()...)
	return *private, pubKey
}
