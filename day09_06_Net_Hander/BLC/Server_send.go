package BLC

import (
	"fmt"
	"net"
	"io"
	"bytes"
	"log"
)

//发送version信息
func sendVersion(toAddr string,bc *Blockchain)  {
	//1.获取当前区块链最新区块高度
	bestHeight:=bc.GetBestHeight()
	//2.创建version对象,第三个参数是要发的是自己节点的地址
	version:=Version{NODE_VERSION,bestHeight,nodeAddress}
	//3.将version序列化
	payload:=gobEncode(version)
	//4.拼接命令+数据
	request:=append(commandToBytes(COMMAND_VERSION),payload...)
	//5.发送数据
	sendData(toAddr,request)
}

//发送消息
func sendData(to string,data []byte)  {
	fmt.Println("当前节点可以发送数据。。。")
	//连接全节点服务器
	conn,err:=net.Dial("tcp",to)
	if err!=nil{
		log.Panic(err)
	}
	defer conn.Close()

	//发送数据
	_,err=io.Copy(conn,bytes.NewReader(data))
	if err!=nil{
		log.Panic(err)
	}
}

