package BLC

import "bytes"

//交易输出
type TXOutput struct {
	//多少钱
	Value int64
	//给谁，这是公钥
	//ScriptPubKey string
	PubKeyHash []byte //公钥哈希
}

//判断TxOutput是否是指定的用户消费
func (txOutput *TXOutput) UnLockWithAddress(address string) bool {
	full_payload:=Base58Decode([]byte(address)) //解码钱包地址
	//取出公钥哈希
	pubKeyHash:=full_payload[1:len(full_payload)-addressCheckSumLen]
	//拿着取出来的公钥哈希和txOutput中的PubKeyHash对比。如果一样，返回真
	return bytes.Compare(pubKeyHash,txOutput.PubKeyHash)==0
}
//根据地址创建一个output对象
func NewTxOutput(value int64,address string) *TXOutput  {
	txOutput:=&TXOutput{value,nil}
	txOutput.Lock(address)
	return txOutput
}
//锁定
func (tx *TXOutput) Lock(address string)  {
	full_payload:=Base58Decode([]byte(address))
	tx.PubKeyHash=full_payload[1:len(full_payload)-addressCheckSumLen]
}