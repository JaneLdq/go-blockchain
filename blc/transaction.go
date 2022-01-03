package blc

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	TxHash []byte
	Vins   []TXInput
	Vouts  []TXOutput
}

type TXInput struct {
	TxHash []byte //交易hash
	Vout   int    //存储TXOutput在Vout里面的索引
	// Signature []byte
	ScriptSig string //用户名
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

//判断当前消费是谁的钱
func (txInput *TXInput) UnLockScriptSigWithAddress(address string) bool {
	return txInput.ScriptSig == address
}

//判断当前消费是谁的钱
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, address}
	// txo.Lock([]byte(address))

	return txo
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {

	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := NewTXOutput(10, to)
	txCoinbase := &Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}

	if data == "Genesis Block" {
		txCoinbase.TxHash = []byte{}
		return txCoinbase
	}
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

// transaction occurs when 'send -f -t -m'
func NewSimpleTransaction(from string, to string, amount int, blc *Blockchain, txs []*Transaction) *Transaction {

	money, spendableUTXODic := blc.FindSpendableUTXOS(from, amount, txs)

	//{hash1:[0], hash2:[2,3]}
	var txInputs []TXInput
	var txOutputs []TXOutput

	for txHash, indexArray := range spendableUTXODic {
		for _, index := range indexArray {
			txHashBytes, _ := hex.DecodeString(txHash)
			txin := TXInput{txHashBytes, index, from}
			txInputs = append(txInputs, txin)
		}
	}

	//expense

	//transaction
	txout := TXOutput{amount, to}
	txOutputs = append(txOutputs, txout)
	//change
	txout = TXOutput{money - amount, from}
	txOutputs = append(txOutputs, txout)
	// utxo
	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	tx.TransactionHash()

	return tx
}

func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}
