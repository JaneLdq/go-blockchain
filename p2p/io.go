package p2p

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"bytes"
	"crypto/elliptic"
	"log"
)

const NODE_FILE_PATH = "node_%d.dat"

func (node *Node) loadFromFile(nodeId uint) error {
	nodeFile := buildNodeFilePath(nodeId)
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

func (node *Node) saveToFile() {
	var content bytes.Buffer
	nodeFile := buildNodeFilePath(node.Port)

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

func buildNodeFilePath(nodeId uint) string {
	return fmt.Sprintf(NODE_FILE_PATH, nodeId)
}