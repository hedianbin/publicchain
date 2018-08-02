package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type Transaction struct {
	//1.交易Hash
	TxHash []byte
	//2.交易输入
	Vins []*TXInput
	//3.交易输出
	Vouts []*TXOutput
}

//创世区块创建时的Transaction
func NewCoinbaseTransaction(address string) *Transaction  {
	//创世区块交易输入,代表消费了
	txInput:=&TXInput{[]byte{},-1,"Genesis Data"}
	//创世区块的交易输出，给address转帐10块，代表未消费
	txOutput:=&TXOutput{10,address}
	//创建Transaction对象
	txCoinbase:=&Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}
	//调用HashTransaction方法设置hash值
	txCoinbase.HashTransaction()
	//返回Transaction对象
	return txCoinbase
}

//序列化Transaction对象并生成hash
func (tx *Transaction) HashTransaction()  {
	//创建bytes.Buffer
	var result bytes.Buffer
	//通过gob.NewEncoder,打包&result
	encoder:=gob.NewEncoder(&result)
	//将Transaction进行序列化
	err:=encoder.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	//将序列化的Transaction生成hash
	hash:=sha256.Sum256(result.Bytes())
	//将生成后的hash传到Transaction的TxHash属性中
	tx.TxHash=hash[:]
}