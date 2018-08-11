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

//钱包文件
const walletFile="wallet_%s.dat"

type Wallets struct {
	//一个字符串对应一个钱包
	Wallets map[string]*Wallet

}

//创建一个钱包，或者抓取已经存在的钱包
func NewWallets(nodeID string) (*Wallets,error) {
	wallets:=Wallets{}
	//开辟内存
	wallets.Wallets=make(map[string]*Wallet)
	err:=wallets.LoadFromFile(nodeID)
	return &wallets,err
}
//创建一个钱包
func (ws *Wallets) CreateWallet() string{
	wallet:=NewWallet()//创建一个钱包
	//创建一个钱包文件的地址
	address:=fmt.Sprintf("%s",wallet.GetAddress())
	ws.Wallets[address]=wallet //把创建的钱包保存
	//返回钱包地址
	return address
}
//抓取所有的钱包的地址
func (ws *Wallets) GetAddresses(nodeID string) []string{
	var addresses []string //装所有的钱包地址
	for address:=range ws.Wallets{ //循环遍历钱包
		addresses=append(addresses,address)
	}
	return addresses //返回所有钱包地址
}
//抓取一个钱包
func (ws *Wallets) GetWallet(address string) Wallet{
	return *ws.Wallets[address]
}


//从文件中读取钱包
func (ws *Wallets) LoadFromFile(nodeID string) error {
	//打印钱包文件名，按nodeID打印钱包文件名，比如nodeID是hedianbin,那么打印出来就是wallet_hedianbin.dat
	mywalletfile:=fmt.Sprintf(walletfile,nodeID)
	//判断钱包文件在不在，不在就返回err
	if _,err:=os.Stat(mywalletfile);os.IsNotExist(err){
		return err
	}
	//读取文件
	fileContent,err:=ioutil.ReadFile(mywalletfile)
	if err!=nil{
		log.Panic(err)
	}
	//读取文件二进制并解析
	var wallets Wallets //钱包集合
	gob.Register(elliptic.P256()) //注册加密解密
	//将读取到的钱包文件进行解码
	decoder:=gob.NewDecoder(bytes.NewReader(fileContent))
	err=decoder.Decode(&wallets) //解码
	if err!=nil{
		log.Panic(err)
	}
	ws.Wallets=wallets.Wallets
	return nil
}

//钱包保存到文件
func (ws *Wallets) SaveToFile(nodeID string) {
	var content bytes.Buffer
	//打印钱包文件名，按nodeID打印钱包文件名，比如nodeID是hedianbin,那么打印出来就是wallet_hedianbin.dat
	mywalletfile:=fmt.Sprintf(walletfile,nodeID)
	gob.Register(elliptic.P256()) //注册加密算法
	encoder:=gob.NewEncoder(&content) //创建编码对象
	err:=encoder.Encode(ws) //将钱包编码
	if err!=nil{
		log.Panic(err)
	}
	//保存钱包文件
	err=ioutil.WriteFile(mywalletfile,content.Bytes(),0644)
	if err!=nil{
		log.Panic(err)
	}
}