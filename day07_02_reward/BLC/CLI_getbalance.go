package BLC

import (
	"fmt"
	"os"
)

//查询余额
func (cli *CLI) GetBalance(address string)  {
	//获取Blockchain对象
	bc:=GetBlockchainObject()
	if bc==nil{
		fmt.Println("没有BlockChain,无法查询。。。")
		os.Exit(1)
	}
	//用完关闭数据库
	defer bc.DB.Close()
	//执行GetBalance方法，查询出余额
	total:=bc.GetBalance(address,[]*Transaction{})
	//输出查到的余额
	fmt.Printf("%s,余额是：%d\n",address,total)
}

