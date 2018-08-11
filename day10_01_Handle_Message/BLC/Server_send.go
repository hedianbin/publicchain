package BLC

import (
	"fmt"
	"net"
	"io"
	"bytes"
	"log"
)

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
//向对方要所有的区块hash
func sendGetBlockHashs(toAddr string)  {
	//创建一个管对方要数据的对象，只传自己的地址过去就行了
	getBlocks:=GetBlocks{nodeAddress}
	//序列化
	payload:=gobEncode(getBlocks)
	//拼接命令+数据
	request:=append(commandToBytes(COMMAND_GETBLOCKHASHS),payload...)
	//发送数据
	fmt.Printf("%s向%s索引所有区块的hash\n",nodeAddress,toAddr)
	sendData(toAddr,request)
}

//发送所有区块hashes
func sendInv(toAddr string,kind string,data [][]byte)  {
	inv:=Inv{nodeAddress,kind,data}
	//序列化要发送的数据
	payload:=gobEncode(inv)
	//拼接要发送的数据，INV就代表告诉对方去处理我发过去的所有hash
	request:=append(commandToBytes(COMMAND_INV),payload...)
	//发送数据
	fmt.Printf("%s已经把拿到的所有区块hash给%s发过去了\n",nodeAddress,toAddr)
	sendData(toAddr,request)
}
//我需要这个hash给对方。对方对应我这个hash给我对应区块数据
func sendGetBlockData(toAddr string,kind string,hash []byte)  {
	getData:=GetData{nodeAddress,kind,hash}
	payload:=gobEncode(getData)
	request:=append(commandToBytes(COMMAND_GETBLOCKDATA),payload...)
	fmt.Printf("%s向%s索要hash对应的区块数据\n",nodeAddress,toAddr)
	sendData(toAddr,request)
}

//发送block给对方
func sendBlock(toAddr string,block *Block)  {
	blockData:=BlockData{nodeAddress,block.Serialize()}
	payload:=gobEncode(blockData)
	request:=append(commandToBytes(COMMAND_BLOCKDATA),payload...)
	sendData(toAddr,request)
}