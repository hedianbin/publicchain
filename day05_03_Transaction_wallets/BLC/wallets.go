package BLC

import (
	"fmt"
)

type Wallets struct {
	//一个字符串对应一个钱包
	WalletMap map[string]*Wallet

}

//创建一个钱包集合,就是给开辟一个钱包集合的空间对象
func NewWallets() *Wallets {
	walletrs:=&Wallets{}
	//开辟内存
	walletrs.WalletMap=make(map[string]*Wallet)
	return walletrs
}
//创建一个钱包
func (ws *Wallets) CreateNewWallet() {
	wallet:=NewWallet()//创建一个钱包集合对象
	//创建一个钱包的地址
	address:=fmt.Sprintf("%s",wallet.GetAddress())
	ws.WalletMap[address]=wallet //把创建的钱包保存
}
