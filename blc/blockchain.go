package blc

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blockTable = "blocks"
const timeFormat string = "2006-01-02 03:04:05 PM"

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

func (blc *Blockchain) MineNewBlock(from string, to string, amount string) {

	fromArray := JSONToArray(from)
	toArray := JSONToArray(to)
	amountArray := JSONToArray(amount)

	fmt.Println(fromArray)
	fmt.Println(toArray)
	fmt.Println(amountArray)

	var txs []*Transaction

	for index, address := range fromArray {

		value, _ := strconv.Atoi(amountArray[index])
		tx := NewSimpleTransaction(address, toArray[index], value, blc, txs)
		txs = append(txs, tx)

	}

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

	// create new block
	block = NewBlock(txs, block.Header.Hash, block.Header.Height+1)
	// save new block to database
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

//return all the unspent transactions
func (bc *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutPuts := make(map[string][]int)

	for _, tx := range txs {
		if !tx.IsCoinbaseTransaction() {
			for _, in := range tx.Vins {
				//是否能够解锁
				if in.UnLockScriptSigWithAddress(address) {
					key := hex.EncodeToString(in.TxHash)
					spentTXOutPuts[key] = append(spentTXOutPuts[key], in.Vout)
				}

			}
		}
	}

	for _, tx := range txs {

	Vouts1:
		for index, out := range tx.Vouts {
			if out.UnLockScriptPubKeyWithAddress(address) {
				if len(spentTXOutPuts) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutPuts {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _, outIndex := range indexArray {
								if index == outIndex {
									isUnSpentUTXO = true
									continue Vouts1
								}

								if !isUnSpentUTXO {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}

							}

						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}
	}

	blockIterator := bc.Iterator()

	for {
		block := blockIterator.Next()

		fmt.Println(block)

		for i := len(block.Txs) - 1; i >= 0; i-- {
			//txHash
			tx := block.Txs[i]
			//Vins

			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					//是否能够解锁
					if in.UnLockScriptSigWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutPuts[key] = append(spentTXOutPuts[key], in.Vout)
					}

				}
			}

			//Vouts
		Vouts:
			for index, out := range tx.Vouts {
				if out.UnLockScriptPubKeyWithAddress(address) {
					if spentTXOutPuts != nil {

						if len(spentTXOutPuts) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutPuts {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue Vouts
									}
								}
							}
							if !isSpentUTXO {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)

							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}
			}

		}

		fmt.Println(spentTXOutPuts)

		var hashInt big.Int
		hashInt.SetBytes(block.Header.PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}

	return unUTXOs
}

// find the spendable utxos
func (blc *Blockchain) FindSpendableUTXOS(from string, amount int, txs []*Transaction) (int, map[string][]int) {

	utxos := blc.UnUTXOs(from, txs)

	spendableUTXO := make(map[string][]int)
	var value int
	for _, utxo := range utxos {
		value = value + utxo.TxOutput.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}

	if value < amount {
		fmt.Printf("%s's fund is not enough\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO

}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string, nodeID uint) *Blockchain {

	dbFile := fmt.Sprintf(dbFile, string(strconv.Itoa(int(nodeID))))

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
func NewBlockchainWithGenesis(nodeID uint) *Blockchain {

	dbFile := fmt.Sprintf(dbFile, string(strconv.Itoa(int(nodeID))))

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

// print blockchain on one node
func PrintChain(nodeId uint) {
	bc := NewBlockchainWithGenesis(nodeId)
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
				fmt.Printf("in.Txid: %x\n", in.TxHash)
				fmt.Printf("n.Vout: %d\n", in.Vout)
				fmt.Printf("in.PubKey: %s\n", in.ScriptSig)
			}
			fmt.Printf("Vouts: \n")
			for _, out := range tx.Vouts {
				// fmt.Println(out)
				fmt.Println(out.Value)
				fmt.Println(out.ScriptPubKey)
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

func (bc *Blockchain) GetBalance(address string) int64 {
	utxos := bc.UnUTXOs(address, []*Transaction{})

	var amount int64

	for _, utxo := range utxos {
		amount = amount + int64(utxo.TxOutput.Value)
	}

	return amount
}
