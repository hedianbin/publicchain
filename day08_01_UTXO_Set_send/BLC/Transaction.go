package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)

type Transaction struct {
	//1.交易Hash
	TxID []byte
	//2.交易输入
	Vins []*TXInput
	//3.交易输出
	Vouts []*TXOutput
}

//创世区块创建时的Transaction
func NewCoinbaseTransaction(address string) *Transaction  {
	//创世区块交易输入,代表消费了
	txInput:=&TXInput{[]byte{},-1,nil,nil}
	//创世区块的交易输出，给address转帐10块，代表未消费
	txOutput:=NewTxOutput(10,address)
	//创建Transaction对象
	txCoinbase:=&Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}
	//调用HashTransaction方法设置hash值
	txCoinbase.SetID()
	//返回Transaction对象
	return txCoinbase
}

//序列化Transaction对象并生成hash
func (tx *Transaction) SetID()  {
	//创建bytes.Buffer
	var result bytes.Buffer
	//通过gob.NewEncoder,打包&result
	encoder:=gob.NewEncoder(&result)
	//将Transaction进行序列化
	err:=encoder.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	//将序列化的Transaction生成hash
	hash:=sha256.Sum256(result.Bytes())
	//将生成后的hash传到Transaction的TxHash属性中
	tx.TxID =hash[:]
}

//2.转帐时产生的Transaction
//需要传3个参数，from源地址，to目标地址，amount要转多少钱,需要返回*Transaction
func NewSimpleTransaction(from string,to string,amount int64,utxoSet *UTXOSet,txs []*Transaction) *Transaction {
	var txInputs []*TXInput
	var txOutputs []*TXOutput

	//获取本次转帐要使用的output，够用就行
	//total,spentableUTXO:=bc.FindSpentableUTXOs(from,amount,txs)
	//从utxoset中拿钱
	total,spentableUTXO:=utxoSet.FindSpentableUTXOs(from,amount,txs)

	//获取钱包的集合
	wallets:=NewWallets()
	//拿到from的钱包
	wallet:=wallets.WalletMap[from]
	//创建多笔交易中的单笔交易的TxInputs
	//遍历拿到的够用的几个output的交易ID和所在交易中的下标
	for txID,indexArray:=range spentableUTXO{
		txIDstr,_:=hex.DecodeString(txID) //将字符串转成字节数组
		//再遍历某一个map中的[]int中所有的下标，取到下标
		for _,index:=range indexArray{
			//创建txInput,第一个参数是花的是哪个交易中的钱，第二个参数是该交易中的哪笔钱的下标，第三个参数是花的谁的钱
			txInput:=&TXInput{txIDstr,index,nil,wallet.PublicKey}
			//记录所有txInput数组
			txInputs=append(txInputs,txInput)
		}
	}

	//给juncheng转帐4块钱
	txOutput:=NewTxOutput(amount,to)
	//记录这笔输出交易到txOutputs数组中
	txOutputs=append(txOutputs,txOutput)
	//给liyuechunc找零6块
	txOutput=NewTxOutput(total-amount,from)
	//记录这笔输出交易到txOutputs数组中
	txOutputs=append(txOutputs,txOutput)
	//创建Transaction对象
	tx:=&Transaction{[]byte{},txInputs,txOutputs}
	//调用HashTransaction方法设置hash值
	tx.SetID()
	//设置签名
	utxoSet.Blockchain.SignTransaction(tx,wallet.PrivateKey,txs)
	//返回Transaction对象
	return tx
}

//判断tx是否是CoinBase交易
func (tx *Transaction) IsCoinBaseTransaction() bool  {
	return len(tx.Vins[0].TxID)==0 && tx.Vins[0].Vout==-1
}

//签名
/*
签名：为了对一笔交易进行签名
私钥：
要获取交易的Input,引用的output,所在的之前的交易
*/
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey,prevTxsmap map[string]*Transaction)  {
	//1.判断当前的tx是否是coinbase交易
	if tx.IsCoinBaseTransaction(){
		return
	}

	//2.拿到刚才传过来的当前的transaction中存储了txid对应的output的map数组
	// 获取当前的txs中的所有的input对应的output所在的tx存不存在，如果不存在，无法进行签名
	for _,input:=range tx.Vins{ //遍历当前交易所有的Vins
		if prevTxsmap[hex.EncodeToString(input.TxID)]==nil{
			log.Panic("当前的input，没有找到对应的output所在的Transaction,无法签名。。")
		}
	}
	//重新构建一份要签名的副本数据
	txCopy:=tx.TrimmedCopy()
	//遍历拿到的副本中的数据进行遍历
	for index,input:=range tx.Vins{
		//从map中拿到当前input对应的tx
		pervTx:=prevTxsmap[hex.EncodeToString(input.TxID)]
		//将副本中的签名置空，双重保险，保证签名一定为空
		txCopy.Vins[index].Signature=nil
		//取出当前input中引用的下标的vout所对应的副本中的vouts中的对应的那笔output中的PubKeyHash
		txCopy.Vins[index].PublicKey=pervTx.Vouts[input.Vout].PubKeyHash
		//产生要签名的数据
		signData:=txCopy.NewTxHash()
		//为了方便下一个input,将PublicKey再置为空
		txCopy.Vins[index].PublicKey=nil

		//开始签名
		/*
		1.第一个参数是随机内存数
		2.第二个参数是参数传过来的私钥
		3.第三个参数是将设置好的txCopy做sha256获取到的hash数据
		*/
		r,s,err:=ecdsa.Sign(rand.Reader,&privateKey,signData)
		if err!=nil{
			log.Panic(err)
		}
		//拼接r+s,就拿到了签名
		sign:=append(r.Bytes(),s.Bytes()...)
		tx.Vins[index].Signature=sign
	}
}

//获取要签名的tx的副本
//要签名的tx中，并不是所有数据都要作为签名数据只是一部分
/*
需要的数据如下
TxID

Inputs
	txid,vout

Outputs
	value,pubkeyhash
注意，除了Inputs中的sign,publickey不要以外，其它都要
*/
//属于tx的方法，处理完副本后返回处理好的tx副本数据
func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs []*TXInput //用于存储所有的TxInput
	var outputs []*TXOutput //用于存储所有的TxOutput
	for _,in:=range tx.Vins{ //遍历所有的input
	//将所有的input的signature和PublicKey置空，然后追加到inputs中，这样就拿到了所有的处理好的TxInput
		inputs=append(inputs,&TXInput{in.TxID,in.Vout,nil,nil})
	}
	for _,out:=range tx.Vouts{ //遍历所有的output
	//全部数据都要
		outputs=append(outputs,&TXOutput{out.Value,out.PubKeyHash})
	}
	//创建新的transaction
	txCopy:=&Transaction{tx.TxID,inputs,outputs}
	return txCopy
}

//序列化
func (tx *Transaction) Serialize() []byte  {
	var buf bytes.Buffer
	encoder:=gob.NewEncoder(&buf)
	err:=encoder.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	return buf.Bytes()
}

//将当前交易生成hash
func (tx *Transaction) NewTxHash() []byte {
	txCopy:=tx            //将当前交易生成一个副本
	txCopy.TxID =[]byte{} //将副本中的txHash置空
	//生成hash
	hash:=sha256.Sum256(txCopy.Serialize())
	//返回hash
	return hash[:]

}

//验证交易
//验证原理：
//用公钥+要签名的数据+r+s，如果为false就说明没验证通过
func (tx *Transaction) Verifity(prevTxs map[string]*Transaction) bool {
	//1.如果是coinbase交易，不需要验证
	if tx.IsCoinBaseTransaction(){
		return true
	}
	//遍历当前传过来的要验证的交易
	for _,input:=range prevTxs{
		if prevTxs[hex.EncodeToString(input.TxID)]==nil{
			log.Panic("当前的input同有找到对应的Transaction,无法验证。。")
		}
	}
	//验证
	txCopy:=tx.TrimmedCopy() //拷贝交易副本

	curev:=elliptic.P256() //创建椭圆曲线加密算法

	//遍历当前交易中的每一笔input数据
	for index,input:=range tx.Vins{
		//原理：再次获取要签名的数据+公钥哈希+签名
		/*
		验证签名的有效性：
		第一个参数：公钥
		第二个参数：签名的数据
		第三，四个参数：r,s
		*/
		//1.获取要签名的数据，主要是用到里面的vouts里面对应的那个output里的公钥哈希
		prevTx:=prevTxs[hex.EncodeToString(input.TxID)]
		//input.Signature=nil
		//input.PublicKey=pervTx.Vouts[input.Vout].PubKeyHash
		//input.TxID =txCopy.NewTxHash() //拿到要签名的数据
		//input.PublicKey=nil            //再次置空
		txCopy.Vins[index].Signature=nil
		txCopy.Vins[index].PublicKey=prevTx.Vouts[input.Vout].PubKeyHash
		txCopy.TxID=txCopy.NewTxHash()
		txCopy.Vins[index].PublicKey=nil

		//获取公钥,将公钥切成x,y,各32bit
		x:=big.Int{}
		y:=big.Int{}
		keyLen:=len(input.PublicKey) //获取input.PublicKey长度值，方便切
		x.SetBytes(input.PublicKey[:keyLen/2])
		y.SetBytes(input.PublicKey[keyLen/2:])

		//再次生成公钥哈希
		rawPublicKey:=ecdsa.PublicKey{curev,&x,&y}

		//获取签名，r,s
		r:=big.Int{}
		s:=big.Int{}
		signLen:=len(input.Signature) //获取签名的长度，方便切
		r.SetBytes(input.Signature[:signLen/2])
		s.SetBytes(input.Signature[signLen/2:])
		//验证签名
		/*
		参数1为新生成的公钥哈希，
		第二个参数为处理好的交易数据副本。
		第三个参数为副本中的交易hash
		第四，第五个参数为用input.Signature拆开的r,s
		 */
		if ecdsa.Verify(&rawPublicKey,txCopy.TxID,&r,&s)==false{
			return false
		}
	}
	return true
}