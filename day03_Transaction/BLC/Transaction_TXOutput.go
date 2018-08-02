package BLC

//交易输出
type TXOutput struct {
	//多少钱
	Value int64
	//给谁，这是公钥
	ScriptPubKey string
}