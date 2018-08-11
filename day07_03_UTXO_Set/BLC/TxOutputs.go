package BLC

import (
	"bytes"
	"encoding/gob"
	"github.com/labstack/gommon/log"
)

//存储所有未花费的UTXO
type TxOutputs struct {
	UTXOs []*UTXO
}

//序列化TxOutputs
func (outs *TxOutputs) Serialize() []byte {
	var buf bytes.Buffer
	encoder:=gob.NewEncoder(&buf)
	err:=encoder.Encode(outs)
	if err!=nil{
		log.Panic(err)
	}
	return buf.Bytes()
}