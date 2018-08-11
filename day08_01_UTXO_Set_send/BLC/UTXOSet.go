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

//从utxoset表中拿要花的钱，够用就停
func (utxoSet *UTXOSet) FindSpentableUTXOs(from string,amount int64,txs []*Transaction) (int64,map[string][]int) {
	var total int64 //存储from的余额
	//用于存储本次转帐需要用的utxo
	spentableUTXOMap := make(map[string][]int)
	//查询未打包的utxo
	unPackageSpentableUTXOs:=utxoSet.FindUnpackeSpentableUTXO(from,txs)
	//遍历获取到的所有的utxo
	for _,utxo:=range unPackageSpentableUTXOs{
		total+=utxo.Output.Value //取出本人所有能花的余额
		txIDStr:=hex.EncodeToString([]byte(utxo.TxID)) //转换[]byte to string
		//将需要转帐的钱存到map中
		spentableUTXOMap[txIDStr]=append(spentableUTXOMap[txIDStr],utxo.Index)
		if total>amount{ //如果转帐的钱够用了，就直接返回拿到的钱
			return total,spentableUTXOMap
		}
	}
	//如查未打包的utxo中的钱不够用
	//查询utxoset表中可用的utxo
	err:=utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(utxosettable)) //查表
		if b!=nil{
			//查询
			c:=b.Cursor()
			//遍历数据库中所有的txoutputs
			dbLoop:
			for k,v:=c.First();k!=nil;k,v=c.Next(){
				txOutputs:=DeserializeTxOutputs(v) //反序列化txOutputs
				for _,utxo:=range txOutputs.UTXOs{ //遍历txOutputs中所有的utxos,拿到其中的所有utxo
					if utxo.Output.UnLockWithAddress(from){ //判断是不是转帐者本人
						total+=utxo.Output.Value //把钱累加
						txIDStr:=hex.EncodeToString(utxo.TxID) //[]byte to string
						//将拿到的钱都放到map中
						spentableUTXOMap[txIDStr]=append(spentableUTXOMap[txIDStr],utxo.Index)
						if total>=amount{ //如果钱够用，跳出最外层循环
							break dbLoop
						}
					}
				}
			}
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	//返回拿到的余额和可以花的钱
	return total,spentableUTXOMap
}

//查询未打包的tx中，可以使用的utxo
func (utxoSet *UTXOSet) FindUnpackeSpentableUTXO(from string,txs []*Transaction) []*UTXO {
	//存储可以使用的未花费的utxo
	var unUTXOs []*UTXO

	//存储已经花费的input
	spentedMap:=make(map[string][]int)

	//倒序遍历每个未打包的交易，去取得每个交易中的未花费的utxo
	for i:=len(txs)-1;i>=0;i--{
		unUTXOs=caculate(txs[i],from,spentedMap,unUTXOs)
	}
	//返回utxos
	return unUTXOs
}