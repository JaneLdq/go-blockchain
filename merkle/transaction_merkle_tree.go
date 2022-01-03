package merkletree

import (
	"go-blockchain/blc"
	"log"
	"reflect"
)

type TransactionContent blc.Transaction
  
func (t TransactionContent) CalculateHash() ([]byte, error) {  
	return t.TxHash, nil
}

func (t TransactionContent) Equals(other Content) (bool, error) {
	return reflect.DeepEqual(t, other), nil	
}

// 使用案例
func test() {
	ti := blc.TXInput{
		TxHash: []byte("test"),
		Vout: 1,
		ScriptSig: "test",
	}
	to := blc.TXOutput{
		Value: 1,
		ScriptPubKey: "test",
	}

	tis := []blc.TXInput{ti}
	tos := []blc.TXOutput{to}

	tx := TransactionContent{
		TxHash: []byte("123123123"),
		Vins: tis,
		Vouts: tos,
	}

	// 创建元数据列表
	var list []Content
	list = append(list, tx)
	list = append(list, tx)
  
	// 根据元数据创建树
	t, err := NewTree(list)
	if err != nil {
	  	log.Fatal(err)
	}
  
	// 获取根hash
	mr := t.MerkleRoot()
	log.Println(mr)
  
	// 校验树的hash
	vt, err := t.VerifyTree()
	if err != nil {
	  	log.Fatal(err)
	}
	log.Println("Verify Tree: ", vt)
  
	// 检查元数据是否在Merkle中
	vc, err := t.VerifyContent(list[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Verify Content: ", vc)
  
	// 获取元数据路径
	hashes, indexes, err := t.GetMerklePath(list[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("path: ", hashes, indexes)


	// 打印树
	log.Println(t.String())
  }


