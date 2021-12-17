package blc

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 16

type PoW struct {
	block  *Block
	target *big.Int
}

func NewPoW(block *Block) *PoW {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &PoW{block, target}

	return pow
}

// Run performs a proof-of-work
func (pow *PoW) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining a new block")
	for {
		hash = pow.CalculateHash(nonce)

		if math.Remainder(float64(nonce), 100000) == 0 {
			fmt.Printf("\r%x", hash)
		}
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *PoW) CalculateHash(nonce int) [32]byte {
	var hash [32]byte
	data := bytes.Join(
		[][]byte{
			pow.block.Header.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Header.Timestamp),
			IntToHex(int64(pow.block.Header.Height)),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	hash = sha256.Sum256(data)
	return hash
}

// Validate validates block's PoW
func (pow *PoW) Validate() bool {
	var hashInt big.Int

	hash := pow.CalculateHash(pow.block.Header.Nonce)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
