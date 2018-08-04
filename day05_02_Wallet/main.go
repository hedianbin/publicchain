package main

import "fmt"
import "publicChain1803/day05_02_Wallet/BLC"

func main() {
	//创建钱包对象，里面有公钥和私钥
	wallet:=BLC.NewWallet()
	address:=wallet.GetAddress()
	fmt.Println(address)
	fmt.Println(string(address))
	fmt.Println(BLC.IsValidAddress(address))
	fmt.Println(BLC.IsValidAddress([]byte("wangergou")))
}
