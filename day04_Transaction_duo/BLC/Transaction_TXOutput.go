package BLC

//交易输出
type TXOutput struct {
	//多少钱
	Value int64
	//给谁，这是公钥
	ScriptPubKey string
}

//判断TxOutput是否是指定的用户消费
func (txOutput *TXOutput) UnLockWithAddress(address string) bool {
	return txOutput.ScriptPubKey==address
}