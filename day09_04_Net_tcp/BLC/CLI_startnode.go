package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) StartNode(nodeID string, mineAddress string) {
	//启动服务器,检验挖矿地址是否为空，或者地址是否有效
	if mineAddress == "" || IsValidAddress([]byte(mineAddress)) {
		//如果没问题。就启动服务器
		startServer(nodeID, mineAddress)
	}else{
		fmt.Println("地址无效。。。")
		os.Exit(1)
	}
}
