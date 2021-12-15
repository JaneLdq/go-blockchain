package dbtest

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	//创建或者打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//创建表 更新表数据
	//err = db.Update(func(tx *bolt.Tx) error {
	//
	//	////创建表BlockBucket
	//	//b, err := tx.CreateBucket([]byte("BlockBucket"))
	//	//if err != nil {
	//	//	return err
	//	//}
	//	//
	//	////往表里存储数据
	//	//if b != nil {
	//	//	err = b.Put([]byte("l"), []byte("Send 100 BTC to cw"))
	//	//	if err != nil {
	//	//		return err
	//	//	}
	//	//}
	//
	//	b := tx.Bucket([]byte("BlockBucket"))
	//
	//	if b != nil {
	//		err = b.Put([]byte("l"), []byte("Send 1030 BTC to re"))
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	//返回nil，以便数据库处理相应操作
	//	return nil
	//})

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("BlockBucket"))

		if b != nil {
			data := b.Get([]byte("l"))
			fmt.Printf("%s", data)
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
