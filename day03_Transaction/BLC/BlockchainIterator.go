package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//迭代器
type BlockchainIterator struct {
	CurrentHash []byte //当前区块的hash值
	DB *bolt.DB //数据库对象
}

//BlockchainIterator返回一个区块
func (blockchainIterator *BlockchainIterator) Next() *Block {
	var block *Block
	err:=blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockTableName))
		//2.如果b!=nil就说明找到了这个表
		if b!=nil{
			//3.获取当前迭代器里的CurrentHash对应的区块数据
			currentBlockBytes:=b.Get([]byte(blockchainIterator.CurrentHash))
			//4.反序列化
			block=DeserializeBlock(currentBlockBytes)
			//更新迭代器里面的CurrentHash值为block.prevBlockHash
			blockchainIterator.CurrentHash=block.PrevBlockHash
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	return block
}
