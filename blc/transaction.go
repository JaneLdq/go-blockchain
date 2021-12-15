package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Transaction struct {
	TxHash []byte
	Vins   []TXInput
	Vouts  []TXOutput
}

type TXInput struct {
	Txid []byte
	Vout int
	// Signature []byte
	PubKey []byte
}

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	// txo.Lock([]byte(address))

	return txo
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to string) *Transaction {

	txin := TXInput{[]byte{}, -1, []byte("Genesis Block")}
	txout := NewTXOutput(110, to)
	txCoinbase := &Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	txCoinbase.TransactionHash()

	return txCoinbase
}

func (tx *Transaction) TransactionHash() {
	txCopy := *tx
	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())
	tx.TxHash = hash[:]

}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
