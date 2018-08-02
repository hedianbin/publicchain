package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)
const targetBit = 16

type ProofOfWork struct {
	//1.要进行验证的区块
	Block *Block
	//2.难度值 big.Int是大数存储
	Target *big.Int
}

//数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int64) []byte  {
	data:=bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		//这里创建一个方法，目的是将[]*Transaction转成字节数组
		pow.Block.HashTransactions(),
		IntToHex(pow.Block.Timestamp),
		IntToHex(int64(nonce)),
		IntToHex(pow.Block.Height),
		IntToHex(targetBit), //这个targetBit可有可无，加不加上都行
	},[]byte{})
	return data
}

//生成有效的hash和nonce值
func (ProofOfWork *ProofOfWork) Run() ([]byte,int64) {
	//设置nonce初始值为0
	var nonce int64 = 0
	//存储我们新生成的hash
	var hashInt big.Int //这里要用指针，因为后面判断的时候需要传进去指针hashInt
	//用来存储生成的有效hash
	var hash [32]byte
	for{
		//1.将Block的属性拼接成字节数组
		dataBytes:=ProofOfWork.prepareData(nonce)
		//2.生成hash
		hash=sha256.Sum256(dataBytes)
		//打印生成hash过程，\r这个意思就是不换行了，覆盖掉上一次输出的值,只显示最新的hash
		//%x，输出十六进制编码
		fmt.Printf("\r%x",hash)
		//将hash存储到hashInt
		hashInt.SetBytes(hash[:])
		//3.判断hashInt是否小于Block里面的target,如果满足条件退出死循环，如果不满足继续循环验证。
		if ProofOfWork.Target.Cmp(&hashInt) == 1{
			break //如果条件成立，跳出循环
		}
		nonce++
	}
	return hash[:],nonce

}
//创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	//1.big.Int对象，初始值为1
	target:=big.NewInt(1)
	//2.左移256-targetBit
	target=target.Lsh(target,256 - targetBit)
	return &ProofOfWork{block,target}
}

//判断hash是否有效
func (ProofOfWork *ProofOfWork) IsValid() bool {
	//1.ProofOfWork.Block.Hash
	//2.ProofOfWork.Target
	//用来存储hash转换的256二进制数
	var hashInt big.Int
	//将hash转换成二进制数256位
	hashInt.SetBytes(ProofOfWork.Block.Hash)
	//判断当前hash有效性
	if ProofOfWork.Target.Cmp(&hashInt)==1{
		return true
	}
	return false
}