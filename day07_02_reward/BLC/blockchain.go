package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"os"
	"strconv"
	"encoding/hex"
	"time"
	"crypto/ecdsa"
	"bytes"
)

//区块链
type Blockchain struct {
	Tip []byte   //最新区块hash值
	DB  *bolt.DB //数据库，用指针
}

//创建迭代器结构体对象
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

//增加区块到区块链里面
/*func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		//1.获取表对象
		b := tx.Bucket([]byte(blockTableName))
		//2.获取数据库中的最新区块字节数组
		blockBytes := b.Get([]byte(blc.Tip))
		//3.将获取到的字节数组反序列化成区块对象
		block := DeserializeBlock(blockBytes)
		if b != nil {
			//4.创建新区块
			newBlock := NewBlock(txs, block.Height+1, block.Hash)
			//5.将区块序列化并且存储到数据库中
			err := b.Put([]byte(newBlock.Hash), []byte(newBlock.Serialize()))
			if err != nil {
				log.Panic(err)
			}
			//6.更新数据库里面"l"对应的hash
			err = b.Put([]byte("l"), []byte(newBlock.Hash))
			if err != nil {
				log.Panic(err)
			}
			//7.更新blockchain中的Tip
			blc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}*/

//1.创建带有创世区块的区块链，返回一个Blockchain结构体，里面装有创世区块
func CreateBlockchainWithGenesisBlock(address string) *Blockchain {
	//判断数据库是否存在
	if DBExists() {
		fmt.Println("创世区块已经存在")
		printUsage()
		os.Exit(1)
	}
	fmt.Println("正在创建创世区块........")
	//创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//用于存储创世区块hash
	var genesisHash []byte
	//更新数据库
	err = db.Update(func(tx *bolt.Tx) error {

		//创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		//如果b==nil说明创建表成功
		if b != nil {
			//创建一个Coinbase
			txCoinbase := NewCoinbaseTransaction(address)
			//1.创建创世区块
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			err = b.Put([]byte(genesisBlock.Hash), []byte(genesisBlock.Serialize()))
			if err != nil {
				log.Panic("数据库存储失败")
			}
			//存储最新的区块hash
			err = b.Put([]byte("l"), []byte(genesisBlock.Hash))
			if err != nil {
				log.Panic("数据库存储失败")
			}
			//拿到创世区块的hash
			genesisHash = genesisBlock.Hash
		}
		return nil
	})
	//返回Blockchain对象
	return &Blockchain{genesisHash, db}
}

//提供一个功能，查询余额
//这个传的txs没用的，实际查询余额不需要Transaction
func (bc *Blockchain) GetBalance(address string, txs []*Transaction) int64 {
	//获取到所有未花费的交易输出数据
	unUTXOs := bc.UnSpent(address, txs)
	//用来存储所有找到的未花费的钱
	var total int64
	//遍历拿到的所有的未花费的钱
	for _, utxo := range unUTXOs {
		//把所有钱加一起
		total += utxo.Output.Value
	}
	return total //返回找到的余额总数
}

//遍历输出所有区块信息
func (blc *Blockchain) Printchain() {
	//生成迭代器
	blockchainIterator := blc.Iterator()
	//2.利用迭代器的Next()循环遍历所有区块数据
	for {
		//每调用一次Next就取一次上一个区块信息，但是不会停止。下面需要判断
		block := blockchainIterator.Next()
		//用来做判断用的，存储block.PrevBlockHash转成big.Int的数据
		//8.输出反序列化后的区块信息
		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		//fmt.Printf("Data:%v\n",block.Txs)
		timeStamp := time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("Timestamp:%s\n", timeStamp)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.Txs {
			fmt.Printf("TxID:%x\n", tx.TxID)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				//已经花掉的钱的交易的hash
				fmt.Printf("\tTxID:%x\n", in.TxID)
				//in.Vout，代表要消费的某一个TXOutput当前数组里面的索引
				fmt.Printf("\tVout:%d\n", in.Vout)
				//签名
				fmt.Printf("\tSignature:%x\n", in.Signature)
				//原始公钥
				fmt.Printf("\tPublicKey:%x\n", in.PublicKey)
			}
			fmt.Println("Vouts:")
			//遍历交易输出信息
			for _, out := range tx.Vouts {
				//花了多少钱
				fmt.Printf("\tValue:%d\n", out.Value)
				//用户名，钱给谁了
				fmt.Printf("\tPubKeyHash:%x\n", out.PubKeyHash)
			}
		}
		fmt.Println()
		var hashInt big.Int
		//将block.PrevBlockHash转成big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		//判断一下如果两边都是256个0就知道找到了创世区块了。退出循环
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

//判断数据库是否存在
func DBExists() bool {
	//err=nil数据库存在，err!=nil不存在
	_, err := os.Stat(dbName)
	//不存在返回false,存在返回true
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetBlockchainObject() *Blockchain {
	//1.创建或打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte //用来存储最新hash
	//查看数据库数据
	err = db.View(func(tx *bolt.Tx) error {
		//获取表对象
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//读取最新区块的hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	return &Blockchain{tip, db}
}

//去挖掘新的区块，需要传进来3个参数
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {
	//1.通过相关算法建立Transaction数组
	var txs []*Transaction
	//建立一笔交易
	//生成多个Transaction
	for i := 0; i < len(from); i++ {
		//将amount转换成int
		value, _ := strconv.ParseInt(amount[i], 10, 64)
		tx := NewSimpleTransaction(from[i], to[i], value, blockchain, txs)
		txs = append(txs, tx)
	}
	//创建新区块之前验证签名有效性
	//遍历每一笔交易签名的有效性
	for _,tx:=range txs{
		//传过去txs，未打包的交易中也要验签
		if blockchain.VerifityTransaction(tx,txs)==false{
			log.Panic("数字签名验证失败....")
		}
	}

	/*
	奖励：reward
	创建一个coinbase交易---tx
	*/
	coinBaseTransaction:=NewCoinbaseTransaction(from[0])
	txs=append(txs,coinBaseTransaction)

	//用来存储拿到的区块
	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		//1.获取表对象b
		b := tx.Bucket([]byte(blockTableName))
		//2.判断表存不存在
		if b != nil {
			//3.拿到最新区块的hash
			hash := b.Get([]byte("l"))
			//4.拿到最新区块的字节数组
			blockBytes := b.Get(hash)
			//5.反序列化字节数组到区块对象
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	//2.挖掘新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	//3.把新区块存储到数据库中
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		//1.获取表对象b
		b := tx.Bucket([]byte(blockTableName))
		//2.判断表存不存在
		if b != nil {
			//3.将最新的block序列化后存储到数据库中
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				log.Panic("数据库存储失败")
			}
			//4.将最新的区块hash存储到数据库key为“l“中
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic("数据库存储失败")
			}
			//更新Tip
			blockchain.Tip = block.Hash
		}
		return nil
	})
}

//查询出所有未花费的output的txid和index
func (bc *Blockchain) FindSpentableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	//1.根据from获取到的所有的utxo
	var total int64
	spentableMap := make(map[string][]int)
	//1.获取所有的未花费的utxo ：10
	utxos := bc.UnSpent(from, txs)
	//2.找即将使用utxo：3个utxo
	//2.遍历utxos，累加余额，判断，是否如果余额，大于等于要要转账的金额
	for _, utxo := range utxos {
		//累加所有可用的余额
		total += utxo.Output.Value
		//把找到的未花费的标记为已花费
		txIDstr := hex.EncodeToString(utxo.TxID)
		spentableMap[txIDstr] = append(spentableMap[txIDstr], utxo.Index)
		if total >= amount {
			break
		}
	}
	if total < amount {
		fmt.Printf("%s,余额不足，无法转账。。\n", from)
		os.Exit(1)
	}

	return total, spentableMap
}

//设计一个方法，用于获取指定用户的所有的未花费TxOutput
func (bc *Blockchain) UnSpent(address string, txs []*Transaction) []*UTXO {
	/*
	1.遍历数据库，获取每个block--->Txs
	2.遍历每个Txs
		Inputs:将数据记录为已经花费
		Outputs:每个output
	*/
	//存储未花费的UTXO
	var unSpentUTXOs []*UTXO
	//存储已经花费的信息
	spentTxOutputMap := make(map[string][]int)
	//第一部分：先查询本次转账，已经产生了的Transanction
	for i := len(txs) - 1; i >= 0; i-- {
		//查询还没有生成区块的Transaction，用的倒序查询，查出来append到unSpentUTXOs中
		unSpentUTXOs = caculate(txs[i], address, spentTxOutputMap, unSpentUTXOs)
	}
	//第二部分，查询数据库里的Transaction
	//获取迭代器
	it := bc.Iterator()
	for {
		//1.获取每个block
		block := it.Next()
		//2.遍历该block的Txs
		for i := len(block.Txs) - 1; i >= 0; i-- {
			//查询数据库中的区块的Transaction，用的还是倒序查询，查出来append到unSpentUTXOs中
			unSpentUTXOs = caculate(block.Txs[i], address, spentTxOutputMap, unSpentUTXOs)
		}

		//3.判断遍历区块退出，说明已经到创世区块了
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}

	}
	//返回查到的所有的未花费的output
	return unSpentUTXOs
}

//查询所有未花费的output
func caculate(tx *Transaction, address string, spentTxOutputMap map[string][]int, unSpentUTXOs []*UTXO) []*UTXO {
	//遍历每个tx：txID，Vins，Vouts
	//遍历所有的TxInput
	if !tx.IsCoinBaseTransaction() { //如果tx不是CoinBase交易就遍历TxInput,否则就不用遍历TxInput
		for _, txInput := range tx.Vins {
			//判断当前的txInput是不是要address这个人花费的，就是看看txInput第三个参数中的ScriptSig是不是和address相同
			//如果相同就记录，不相同就执行下面语句
			full_payload := Base58Decode([]byte(address))
			pubKeyHash := full_payload[1 : len(full_payload)-addressCheckSumLen]
			if txInput.UnLockWithAddress(pubKeyHash) {
				//txInput的解锁脚本(用户名) 如果和要查询的余额的用户名相同，
				key := hex.EncodeToString(txInput.TxID)
				//将查询到的已花费的这笔钱存到map中，
				spentTxOutputMap[key] = append(spentTxOutputMap[key], txInput.Vout)
				/*
				map[key]-->value
				map[key] -->[]int
				 */
			}
		}
	}
	//遍历所有的TxOutput交易
outputs:
	for index, txOutput := range tx.Vouts {
		//判断遍历的txOutput是否是和要查询余额的这个人的，如果不是就不用执行下面语句
		if txOutput.UnLockWithAddress(address) {
			if len(spentTxOutputMap) != 0 { //如果map记录了txInput，就去过滤
				var isSpentOutput bool //记录是否已花费
				//遍历map
				for txID, indexArray := range spentTxOutputMap {
					//遍历，记录已经花费的下标的数组
					for _, i := range indexArray {
						//判断，如果map中查到的下标记录和当前TxOutput的下标相同，
						//并且，当前区块的交易hash和map中的key相同，就说明这笔钱已经被花费掉了。
						if i == index && hex.EncodeToString(tx.TxID) == txID {
							isSpentOutput = true //标记当前的txOutput已经花费掉了
							continue outputs     //当前区块的tx.Vouts数组就不用遍历了
						}
					}
				}
				//如果此txOutput没有在map中查到被花费掉，就记录一下
				if !isSpentOutput {
					//unSpentTxOutput=append(unSpentTxOutput,txOutput)
					utxo := &UTXO{tx.TxID, index, txOutput}
					unSpentUTXOs = append(unSpentUTXOs, utxo)
				}
			} else {
				//如果map长度为0，证明还没有花费记录，output无需判断
				//unSpentTxOutput=append(unSpentTxOutput,txOutput)
				utxo := &UTXO{tx.TxID, index, txOutput}
				unSpentUTXOs = append(unSpentUTXOs, utxo)
			}
		}
	}
	return unSpentUTXOs
}

//签名
func (bc *Blockchain) SignTransaction(tx *Transaction, prevateKey ecdsa.PrivateKey,txs []*Transaction) {
	//1.判断要签名的tx,如果是coinbase交易直接返回
	if tx.IsCoinBaseTransaction() {
		return
	}
	//2.获取该tx中的Input,引用之前的transaction中的未花费的output
	prevTxs := make(map[string]*Transaction)
	for _, input := range tx.Vins {
		txIDStr := hex.EncodeToString(input.TxID)
		//根据tx中的每个input对应txID，获取要转账的的output所在的Transaction，存入到map[input.TxID]-->Tx
		//将txs传进去的意思是，多笔交易的时候，先从未打包的txs中拿transaction
		prevTxs[txIDStr] = bc.FindTransactionByTxID(input.TxID,txs)
	}
	//3。签名
	tx.Sign(prevateKey, prevTxs)
}

//根据交易ID，获取对应的交易对象
func (bc *Blockchain) FindTransactionByTxID(txID []byte,txs []*Transaction) *Transaction {
	//1.先查找未打包的txs
	//如果在未打包的txs中找到了transaction,下面查库就不需要了
	for _,tx:=range txs{
		if bytes.Compare(tx.TxID,txID)==0{
			return tx
		}
	}
	//遍历数据库，获取block--->Transaction
	iterator := bc.Iterator()
	for {
		//遍历每个区块
		block := iterator.Next()
		for _, tx := range block.Txs { //再遍历每个区块中的transaction
			if bytes.Compare(tx.TxID, txID) == 0 {
				return tx
			}
		}
		//判断结束循环
		var bigInt big.Int
		bigInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&bigInt)==0{
			break
		}
	}
	return &Transaction{}
}

//验证交易的数字签名
//未打包的交易中也要验，所以传过去txs
func (bc *Blockchain) VerifityTransaction(tx *Transaction,txs []*Transaction) bool {
	//要想验证数字签名：需要私钥+数据（tx的副本+之前的交易）
	prevTxs:=make(map[string]*Transaction)
	//遍历要验证的tx交易中的所有的Vins中的所有的input
	for _,input:=range tx.Vins{
		//拿到对应的input.Txid对应的数据库中对应的transaction
		//如果未打包中的txs有，就先从txs中拿
		prevTx:=bc.FindTransactionByTxID(input.TxID,txs)
		//用当前的input.Txid作为key,拿到的transaction作为value，存储到map中
		prevTxs[hex.EncodeToString(input.TxID)]=prevTx
	}
	//验证
	//用上面获取到的map去验证有效性
	return tx.Verifity(prevTxs)
}