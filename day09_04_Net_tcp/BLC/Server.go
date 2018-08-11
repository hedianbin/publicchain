package BLC

import (
	"fmt"
	"net"
	"log"
	"io"
	"bytes"
	"io/ioutil"
)

//全节点地址
var knowNodes=[]string{"localhost:3000"}

//节点服务端启动，可以接收其他节点传来的数据
func startServer(nodeID string,mineAddress string)  {
	//当前节点自己的地址
	nodeAddress:=fmt.Sprintf("localhost:%s",nodeID)
	//监听地址
	listener,err:=net.Listen("tcp",nodeAddress)
	if err!=nil{
		log.Panic(err)
	}
	//延迟关闭连接
	defer listener.Close()

	//如果不是全节点，比如是钱包节点或矿工节点，就给全节点发送一个消息
	if nodeAddress!=knowNodes[0]{
		//是钱包节点，或矿工节点，那么给全节点发一个消息
		sendMessage(knowNodes[0],"我是王二狗，我的地址是："+nodeAddress)
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

func sendMessage(to,data string)  {
	fmt.Println("当前节点可以发送数据。。。")
	//连接全节点服务器
	conn,err:=net.Dial("tcp",to)
	if err!=nil{
		log.Panic(err)
	}
	defer conn.Close()

	//发送数据
	_,err=io.Copy(conn,bytes.NewReader([]byte(data)))
	if err!=nil{
		log.Panic(err)
	}
}
