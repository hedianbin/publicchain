package BLC

//交易输入
type TXInput struct {
	//1.交易的ID
	Txid []byte
	//2.用来存储TXOutput在Transaction结构体的Vouts里面的索引
	Vout int
	//3.用户名，这是数字签名
	ScriptSig string
}
//判断TxInput是否是指定的用户消费
func (txInput *TXInput) UnLockWithAddress(address string) bool {
	return txInput.ScriptSig==address
}