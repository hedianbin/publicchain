package BLC

import (
	"fmt"
	"os"
)

//打印输出数据库中所有区块信息
func (cli *CLI) printchain(nodeID string) {
	DBName:=fmt.Sprintf(dbName,nodeID)
	if DBExists(DBName) == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain := GetBlockchainObject(nodeID)
	defer blockchain.DB.Close()
	blockchain.Printchain(nodeID)
}