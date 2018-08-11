package BLC

import (
	"fmt"
	"os"
)

//打印输出数据库中所有区块信息
func (cli *CLI) printchain() {
	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain := GetBlockchainObject()
	defer blockchain.DB.Close()
	blockchain.Printchain()
}