package BLC

type UTXO struct {
	TxID []byte
	Index int
	Output *TXOutput
}