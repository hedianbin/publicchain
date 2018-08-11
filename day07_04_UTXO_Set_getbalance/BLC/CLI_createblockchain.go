package BLC

import (
	"fmt"
	"os"
)

//创建创世区块
func (cli *CLI) createGenesisBlockchain(address string) {
	blockchain:=CreateBlockchainWithGenesisBlock(address)
	defer blockchain.DB.Close()

	if blockchain==nil{
		fmt.Println("没有数据库。。。")
		os.Exit(1)
	}
	utxoSet:=&UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()
}