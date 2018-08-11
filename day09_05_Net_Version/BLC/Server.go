package BLC

import (
	"fmt"
	"net"
	"log"
	"io"
	"bytes"
	"io/ioutil"
)



//节点服务端启动，可以接收其他节点传来的数据
func startServer(nodeID string,mineAddress string)  {
	//当前节点自己的地址
	nodeAddress=fmt.Sprintf("localhost:%s",nodeID)
	//监听地址
	listener,err:=net.Listen("tcp",nodeAddress)
	if err!=nil{
		log.Panic(err)
	}
	//延迟关闭连接
	defer listener.Close()

	//拿到bc对象
	bc:=GetBlockchainObject(nodeID)
	defer bc.DB.Close()
	//如果不是全节点，比如是钱包节点或矿工节点，就给全节点发送一个消息
	if nodeAddress!=knowNodes[0]{
		//给全节点发送版本信息
		sendVersion(knowNodes[0],bc)
	}
	//如果是全节点。等待客户端连入
	for{
		//等待客户端连入，阻塞了
		conn,err:=listener.Accept()
		if err!=nil{
			log.Panic(err)
		}
		fmt.Println("发送方已经连入：",conn.RemoteAddr())
		//读取数据
		request,err:=ioutil.ReadAll(conn)
		if err!=nil{
			log.Panic(err)
		}
		fmt.Printf("接收到的数据是：%s\n",request)
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
