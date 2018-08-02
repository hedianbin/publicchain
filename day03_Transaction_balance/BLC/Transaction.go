package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
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
	txCoinbase.SetID()
	//返回Transaction对象
	return txCoinbase
}

//序列化Transaction对象并生成hash
func (tx *Transaction) SetID()  {
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

//2.转帐时产生的Transaction
//需要传3个参数，from源地址，to目标地址，amount要转多少钱,需要返回*Transaction
func NewSimpleTransaction(from string,to string,amount int64) *Transaction {
	var txInputs []*TXInput
	var txOutputs []*TXOutput
	//代表消费掉了liyuechun的10块钱
	b,_:=hex.DecodeString("6584e2f3e833a7ced439438f1d79b28b2d626b67e8a179d0f643a0463ea5efb9")
	txInput:=&TXInput{b,1,from}
	//所有已消费的钱
	txInputs=append(txInputs,txInput)
	//给juncheng转帐4块钱
	txOutput:=&TXOutput{amount,to}
	//记录这笔输出交易到txOutputs数组中
	txOutputs=append(txOutputs,txOutput)
	//给liyuechunc找零6块
	txOutput=&TXOutput{6-amount,from}
	//记录这笔输出交易到txOutputs数组中
	txOutputs=append(txOutputs,txOutput)
	//创建Transaction对象
	tx:=&Transaction{[]byte{},txInputs,txOutputs}
	//调用HashTransaction方法设置hash值
	tx.SetID()
	//返回Transaction对象
	return tx
}

//判断tx是否是CoinBase交易
func (tx *Transaction) IsCoinBaseTransaction() bool  {
	return len(tx.Vins[0].Txid)==0 && tx.Vins[0].Vout==-1
}