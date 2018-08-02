package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"time"
	"os"
)
//数据库名字
const dbName = "blockchain.db"
//表的名字
const blockTableName = "blocks"

//区块链
type Blockchain struct {
	Tip []byte //最新区块hash值
	DB *bolt.DB //数据库，用指针
}


//创建迭代器结构体对象
func (blockchain *Blockchain) Iterator() *BlockchainIterator  {
	return &BlockchainIterator{blockchain.Tip,blockchain.DB}
}


//增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {
	err:=blc.DB.Update(func(tx *bolt.Tx) error {
		//1.获取表对象
		b:=tx.Bucket([]byte(blockTableName))
		//2.获取数据库中的最新区块字节数组
		blockBytes:=b.Get([]byte(blc.Tip))
		//3.将获取到的字节数组反序列化成区块对象
		block:=DeserializeBlock(blockBytes)
		if b!=nil{
			//4.创建新区块
			newBlock:=NewBlock(txs,block.Height+1,block.Hash)
			//5.将区块序列化并且存储到数据库中
			err:=b.Put([]byte(newBlock.Hash),[]byte(newBlock.Serialize()))
			if err!=nil{
				log.Panic(err)
			}
			//6.更新数据库里面"l"对应的hash
			err=b.Put([]byte("l"),[]byte(newBlock.Hash))
			if err!=nil{
				log.Panic(err)
			}
			//7.更新blockchain中的Tip
			blc.Tip=newBlock.Hash
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
}

//1.创建带有创世区块的区块链，返回一个Blockchain结构体，里面装有创世区块
func CreateBlockchainWithGenesisBlock(address string)   {
	//判断数据库是否存在
	if DBExists(){
		fmt.Println("创世区块已经存在")
		printUsage()
		os.Exit(1)
	}
	fmt.Println("正在创建创世区块........")
	//创建或打开数据库
	db,err:=bolt.Open(dbName,0600,nil)
	if err!=nil{
		log.Fatal(err)
	}
	defer db.Close()
	//更新数据库
	err=db.Update(func(tx *bolt.Tx) error {

		//创建数据库表
		b,err:=tx.CreateBucket([]byte(blockTableName))
		if err!=nil{
			log.Panic(err)
		}

		//如果b==nil说明创建表成功
		if b!=nil{
			//创建一个Coinbase
			txCoinbase:=NewCoinbaseTransaction(address)
			//1.创建创世区块
			genesisBlock:=CreateGenesisBlock([]*Transaction{txCoinbase})
			err=b.Put([]byte(genesisBlock.Hash),[]byte(genesisBlock.Serialize()))
			if err!=nil{
				log.Panic("数据库存储失败")
			}
			//存储最新的区块hash
			err=b.Put([]byte("l"),[]byte(genesisBlock.Hash))
			if err!=nil{
				log.Panic("数据库存储失败")
			}

		}
		return nil
	})


}
//遍历输出所有区块信息
func (blc *Blockchain) Printchain()  {
	//生成迭代器
	blockchainIterator:=blc.Iterator()
	//2.利用迭代器的Next()循环遍历所有区块数据
	for{
		//每调用一次Next就取一次上一个区块信息，但是不会停止。下面需要判断
		block:=blockchainIterator.Next()
		//用来做判断用的，存储block.PrevBlockHash转成big.Int的数据
		//8.输出反序列化后的区块信息
		fmt.Printf("Height:%d\n",block.Height)
		fmt.Printf("PrevBlockHash:%x\n",block.PrevBlockHash)
		//fmt.Printf("Data:%v\n",block.Txs)
		timeStamp:=time.Unix(block.Timestamp,0).Format("2006-01-02 15:04:05")
		fmt.Printf("Timestamp:%s\n",timeStamp)
		fmt.Printf("Hash:%x\n",block.Hash)
		fmt.Printf("Nonce:%d\n",block.Nonce)
		for _,tx:=range block.Txs{
			fmt.Println(tx)
		}
		fmt.Println()
		var hashInt big.Int
		//将block.PrevBlockHash转成big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		//判断一下如果两边都是256个0就知道找到了创世区块了。退出循环
		if big.NewInt(0).Cmp(&hashInt)==0{
			break
		}
	}
}

//判断数据库是否存在
func DBExists() bool {
	//err=nil数据库存在，err!=nil不存在
	_,err:=os.Stat(dbName)
	//不存在返回false,存在返回true
	if os.IsNotExist(err){
		return false
	}
	return true
}

func GetBlockchainObject() *Blockchain  {
	//1.创建或打开数据库
	db,err:=bolt.Open(dbName,0600,nil)
	if err!=nil{
		log.Fatal(err)
	}
	var tip []byte //用来存储最新hash
	//查看数据库数据
	err = db.View(func(tx *bolt.Tx) error {
		//获取表对象
		b:=tx.Bucket([]byte(blockTableName))
		if b!=nil{
			//读取最新区块的hash
			tip=b.Get([]byte("l"))
		}
		return nil
	})
	return &Blockchain{tip,db}
}