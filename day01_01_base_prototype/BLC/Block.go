package BLC


import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
	"fmt"
	"encoding/gob"
	"log"
)

//定义区块
type Block struct {
	//1.区块高度，编号，第几个区块
	Height int64
	//2.上一个区块的PreHash
	PrevBlockHash []byte
	//3.交易数据Data,
	//不管是转帐，部署合约，实际上交易数据最终转换成transaction,就会把一些东西装进来。
	Data []byte
	//4.时间戳
	Timestamp int64
	//5.当前Hash
	Hash []byte
	//6.Nonce
	Nonce int64
}
//生成当前区块hash
func (block *Block) SetHash()  {
	//1.将height转化为字节数组
	heightBytes:=IntToHex(block.Height)
	//2.timestamp转化为字节数组
	//(1)将int64转化为字符串
	//第二个参数的范围为2~36,代表将时间戳转换为多少进制的字符串
	timeString:=strconv.FormatInt(block.Timestamp,2)
	//（2）将字符串转字节数组
	timeBytes:=[]byte(timeString)
	//3.拼接所有属性
	blockBytes:=bytes.Join([][]byte{heightBytes,block.PrevBlockHash,block.Data,timeBytes,block.Hash},[]byte{})
	//4.将拼接成的字节数组生成hash
	hash:=sha256.Sum256(blockBytes)
	block.Hash=hash[:]
}

//1.创建新的区块
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {
	timestamp := time.Now().Unix()
	//创建区块
	block := &Block{height, prevBlockHash, []byte(data), timestamp, nil,0}

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
func CreateGenesisBlock(data string) *Block  {
	return NewBlock(data,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
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