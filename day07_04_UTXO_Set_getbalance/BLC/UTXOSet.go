package BLC

import (
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
	"encoding/hex"
)

/*
将UTXO持久化
数据库：blockchain.db
表1：blocks 存储所有区块
表2：utxoset 存储所有utxo
查询余额，转帐
*/
//拿到blockchain对象，db,tip
type UTXOSet struct {
	Blockchain *Blockchain
}

//utxoset表名
const utxosettable = "utxoset"

//将UTXO持久化
func (utxoset *UTXOSet) ResetUTXOSet()  {
	//查询block块中的所有的未花费utxo
	err:=utxoset.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		//1.如果utxoset表存在，删除
		b:=tx.Bucket([]byte(utxosettable))
		if b!=nil{
			//删表
			err:=tx.DeleteBucket([]byte(utxosettable))
			if err!=nil{
				log.Panic("重置时，删除表失败")
			}
		}
		//2.创建utxoset
		b,err:=tx.CreateBucket([]byte(utxosettable))
		if err!=nil{
			log.Panic("重置时，创建表失败")
		}
		//如果表创建成功
		if b!=nil{
			//去到库中查到所有的未花费的utxo存到map中
			unUTXOMap:=utxoset.Blockchain.FindUnspentUTXOMap()

			//遍历拿到的所有未花费的txOutputs
			for txIDStr,outs:=range unUTXOMap{
				txID,_:=hex.DecodeString(txIDStr) //字符串转成[]byte
				//将txOutputs序列化后存储到表中
				b.Put(txID,outs.Serialize())
			}
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
}

//从到utxoset中去查询余额
func (utxoSet *UTXOSet) GetBalance(address string) int64 {
	//去utxoset中查询所有未花费的utxo
	utxos:=utxoSet.FindUnspentUTXOsByAddress(address)
	var total int64 //用来记录余额
	for _,utxo:=range utxos{
		//累加所有的未花费的金额
		total+=utxo.Output.Value
	}
	//返回找到的金额
	return total
}

//根据地址查到utxoset表中所有的utxo
func (utxoSet *UTXOSet) FindUnspentUTXOsByAddress(address string) []*UTXO {
	//存储查到的所有的未花费的utxo
	var utxos []*UTXO
	err:=utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(utxosettable)) //打开utxoset表
		if b!=nil{
			//获取表中的所有的utxo
			c:=b.Cursor()
			//遍历数据库,拿到对应address的所有的txInputs
			for k,v:=c.First();k!=nil;k,v=c.Next(){
				//将每一个txInoutputs反序列化
				txOutputs:=DeserializeTxOutputs(v)
				//遍历反序列化后的所有的utxo
				for _,utxo:=range txOutputs.UTXOs{
					//判断是否本人查询
					if utxo.Output.UnLockWithAddress(address){
						//如果是本人查询就把查到的utxos返回
						utxos=append(utxos,utxo)
					}
				}
			}
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	return utxos
}