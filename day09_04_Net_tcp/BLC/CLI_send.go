package BLC

import (
	"fmt"
	"os"
)

//转帐
func (cli *CLI) Send(from []string,to []string,amount []string,nodeID string)  {

	DBName:=fmt.Sprintf(dbName,nodeID)

	if !DBExists(DBName){
		fmt.Println("数据库不存在......")
		os.Exit(1)
	}
	//拿到了带有最新区块hash和db对象的Blockchain对象
	bc:=GetBlockchainObject(nodeID)
	//关闭数据库连接
	defer bc.DB.Close()
	//挖掘新区块，带Transaction交易的
	bc.MineNewBlock(from,to,amount,nodeID)

	//添加更新
	utxoSet :=&UTXOSet{bc}
	utxoSet.Update()
}
