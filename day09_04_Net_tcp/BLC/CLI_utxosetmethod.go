package BLC

import (
	"fmt"
	"os"
)

//重置UTXO，测试代码
func (cli *CLI) UtxoSetMethod(nodeID string)  {
	//获取blockchain对象，db,tip
	bc:=GetBlockchainObject(nodeID)
	if bc==nil{
		fmt.Println("没有数据库，无法获取utxo....")
		os.Exit(1)
	}
	defer bc.DB.Close()
	//拿到所有未花费的utxo
	unSpentUTXOsMap:=bc.FindUnspentUTXOMap()
	fmt.Println("长度：",len(unSpentUTXOsMap))
	//遍历拿到的所有TxOutputs
	for txIDStr,txOutputs:=range unSpentUTXOsMap{
		fmt.Println("交易ID：",txIDStr)
		//拿到某个TxOutputs里的所有的UTXOs中的每个utxo
		for _,utxo:=range txOutputs.UTXOs{
			fmt.Println("\t金额：",utxo.Output.Value)
			fmt.Printf("\t地址：%s\n",GetAddressByPubKeyHash(utxo.Output.PubKeyHash))
		}
		fmt.Println("===================================================================")
	}
	//重置utxoset表
	utxoSet:=&UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}
