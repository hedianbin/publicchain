package BLC

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"log"
	"io/ioutil"
	"os"
)

const walletsfile="wallets_%s.dat"

type Wallets struct {
	//一个字符串对应一个钱包
	WalletMap map[string]*Wallet
}

//创建一个钱包集合,就是给开辟一个钱包集合的空间对象
func NewWallets(nodeID string) *Wallets {
	//设置钱包文件名
	walletsFile:=fmt.Sprintf(walletsfile,nodeID)

	//1.钱包文件不存在
	if _,err:=os.Stat(walletsFile);os.IsNotExist(err){
		fmt.Println("钱包文件不存在")
		wallets:=&Wallets{}
		wallets.WalletMap=make(map[string]*Wallet)
		return wallets
	}
	//2.钱包文件存在，读取本地文件，将钱包数据转成钱包对象
	//读取本地的钱包文件中的数据,并反序列化得到钱包集合对象
	wsBytes,err:=ioutil.ReadFile(walletsFile)
	if err!=nil{
		log.Panic(err)
	}
	//将数据变成钱包集合对象
	gob.Register(elliptic.P256())//创建加密算法
	var wallets Wallets //创建钱包集合对象
	//对钱包里的进行反序列化
	reader:=bytes.NewReader(wsBytes)
	decoder:=gob.NewDecoder(reader)
	err=decoder.Decode(&wallets)
	if err!=nil{
		log.Panic(err)
	}
	//返回解码后的钱包集合
	return &wallets

}

//创建钱包
func (ws *Wallets) CreateNewWallet(nodeID string)  {
	//创建钱包集合，或从钱包文件中反序列化
	wallet:=NewWallet()
	//生成一个钱包地址
	address:=wallet.GetAddress()
	fmt.Printf("创建的钱包的地址：%s\n",address)
	//将新生成的钱包地址保存到钱包集合中
	ws.WalletMap[string(address)]=wallet
	ws.saveFile(nodeID)
}

//将钱包集合保存到本地文件中
func (ws *Wallets) saveFile(nodeID string)  {
	walletsFile:=fmt.Sprintf(walletsfile,nodeID)
	//1.将ws对象的数据转成[]byte
	var buf bytes.Buffer
	//序列化钱包集合，序列化过程中，如果被序列化的对象中包含了接口，那么接口需要注册
	gob.Register(elliptic.P256())
	encoder:=gob.NewEncoder(&buf) //注册buf
	err:=encoder.Encode(ws) //序列化钱包集合
	if err!=nil{
		log.Panic(err)
	}
	//拿到序列化后的钱包集合
	wsBytes:=buf.Bytes()
	//2.将数据存储到文件中
	err=ioutil.WriteFile(walletsFile,wsBytes,0644)
	if err!=nil{
		log.Panic(err)
	}
}