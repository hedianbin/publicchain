package BLC

import "bytes"

//交易输入
type TXInput struct {
	//1.交易的ID
	TxID []byte
	//2.用来存储TXOutput在Transaction结构体的Vouts里面的索引
	Vout int
	//3.用户名，这是数字签名
	//ScriptSig string
	//3.解锁脚本
	Signature []byte //数字签名
	PublicKey []byte //原始公钥，钱包里的公钥

}
//判断TxInput是否是指定的用户消费
func (txInput *TXInput) UnLockWithAddress(pubKeyHash []byte) bool {
	//用原始公钥生成公钥哈希
	pubKeyHash2:=PubKeyHash(txInput.PublicKey)
	//用生成的公钥哈希和传过来的对比一下，一样就返回true
	return bytes.Compare(pubKeyHash,pubKeyHash2)==0
}