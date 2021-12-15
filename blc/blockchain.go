package blc

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_db"
const blockTable = "blocks"

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTable))
		blockBytes := b.Get(i.currentHash)
		block = DeserializeBlock(blockBytes)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.Header.PrevBlockHash

	return block
}

func (blc *Blockchain) MineNewBlock(from []string, to []string, amount []string) {
	// fmt.Println(from)
	// fmt.Println(to)
	// fmt.Println(amount)

	var txs []*Transaction

	var block *Block
	blc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTable))
		if b != nil {
			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = DeserializeBlock(blockBytes)
		}

		return nil
	})

	block = NewBlock(txs, block.Header.Hash, block.Header.Height+1)

	blc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTable))
		if b != nil {
			b.Put(block.Header.Hash, block.Serialize())
			b.Put([]byte("l"), block.Header.Hash)
			blc.Tip = block.Header.Hash
		}
		return nil
	})

}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {

	if DBExists(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	coinbaseTX := NewCoinbaseTX(address)
	genesis := NewGenesisBlock(coinbaseTX)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockTable))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Header.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Header.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Header.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// creates a new Blockchain with genesis Block
func NewBlockchainWithGenesis() *Blockchain {

	if !DBExists(dbFile) {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTable))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// Iterator returns a BlockchainIterat
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip, bc.DB}

	return bci
}

func DBExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
