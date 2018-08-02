package main

import (
	"publicChain1803/day01_01_base_prototype/BLC"
)

func main() {
	//创世区块
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()
	//新区块

	//添加4个新区块
	blockchain.AddBlockToBlockchain("Send 100RMB To zhangqiang")

	blockchain.AddBlockToBlockchain("Send 200RMB To changjingkong")

	blockchain.AddBlockToBlockchain("Send 300RMB To jiatengying")

	blockchain.AddBlockToBlockchain("Send 50RMB To hedianbin")

	blockchain.Printchain()
}