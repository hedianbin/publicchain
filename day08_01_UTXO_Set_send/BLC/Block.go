package BLC


import (
	"bytes"
	"time"
	"fmt"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

//定义区块
type Block struct {
	//1.区块高度，编号，第几个区块
	Height int64
	//2.上一个区块的PreHash
	PrevBlockHash []byte
	//3.交易数据Data,
	//不管是转帐，部署合约，实际上交易数据最终转换成transaction,就会把一些东西装进来。
	Txs []*Transaction
	//4.时间戳
	Timestamp int64
	//5.当前Hash
	Hash []byte
	//6.Nonce
	Nonce int64
}

//1.创建新的区块
func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {
	timestamp := time.Now().Unix()
	//创建区块
	block := &Block{height, prevBlockHash, txs, timestamp, nil,0}

	//创建工作量证明对象，传入一个block,然后计算出target值，就是前面有几个0的难度值。然后返回ProofOfWork，里面装有block和计算好的target
	pow:=NewProofOfWork(block) //创建pow对象

	//调用工作量证明的Run方法执行挖矿验证,并且返回有效的Hash和Nonce值
	//运行一次就计算一次，计算合法的hash值后，返回有效的hash和nonce
	//比如设定了hash前面000000，如果生成的hash值前面带000000,就算合法的hash,并将对应的nonce和hash返回回来
	hash,nonce:=pow.Run()

	block.Hash=hash
	block.Nonce=nonce
	//添加换行，解决缓冲区问题
	fmt.Println()
	return block
}

//2.单独写一个方法，生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block  {
	return NewBlock(txs,1,make([]byte,32,32))
}

//将区块序列化成字节数组
func (block *Block) Serialize() []byte  {
	//创建bytes.Buffer
	var result bytes.Buffer
	//通过gob.NewEncoder,打包&result
	encoder:=gob.NewEncoder(&result)
	//将区块进行序列化
	err:=encoder.Encode(block)
	if err!=nil{
		log.Panic(err)
	}
	//调用result.Bytes{},返回序列化后的字节数组
	return result.Bytes()
}

//反序列化成区块对象
func DeserializeBlock(blockBytes []byte) *Block  {
	//用来存储反序列化的区块对象
	var block Block
	//解包需要反序列化的字节数组
	decoder:=gob.NewDecoder(bytes.NewReader(blockBytes))
	//调用decoder.Decode将字节数组反序列化后存储到&block中
	err:=decoder.Decode(&block)
	if err!=nil{
		log.Panic(err)
	}
	//返回反序列化后的区块对象
	return &block
}

//将[]*Transaction转换成字节数组
func (block *Block) HashTransactions() []byte {
	//用来存储append后的hash值，是将[]byte通过append后，变成一个二维数组[][]byte
	var txHashes [][]byte
	//用来存储将拼接后的txHashes加密为新的hash数组
	var txHash [32]byte
	//遍历Tsx中的TxHash值,并拼接起来
	for _,tx:=range block.Txs{
		//拼接TxHash
		txHashes=append(txHashes,tx.TxID)
	}
	//将拼接的数据进行sha256加密
	txHash=sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	//返回加密后的hash，需要切一下
	return txHash[:]
}