package blc

type UTXO struct {
	TxHash   []byte //Transaction hash
	Index    int    //Transaction index
	TxOutput TXOutput
}
