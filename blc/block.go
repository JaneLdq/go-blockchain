package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Header *BlockHeader
	Txs    []*Transaction
}

type BlockHeader struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
}

// NewBlock creates and returns Block
func NewBlock(txs []*Transaction, prevBlockHash []byte, height int) *Block {
	blockHeader := &BlockHeader{time.Now().Unix(), prevBlockHash, []byte{}, 0, height}
	block := &Block{blockHeader, txs}
	pow := NewPoW(block)
	nonce, hash := pow.Run()

	block.Header.Hash = hash[:]
	block.Header.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// append all transactions and hash them
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte
	var txHash [32]byte

	for _, tx := range b.Txs {
		transactions = append(transactions, tx.TxHash)
	}
	// mTree := NewMerkleTree(transactions)

	// return mTree.RootNode.Data
	txHash = sha256.Sum256(bytes.Join(transactions, []byte{}))

	return txHash[:]
}
