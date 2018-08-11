package main

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
)

func main() {
	fmt.Println("服务端的程序")
	listen,err:=net.Listen("tcp",":9528")
	if err!=nil{
		log.Panic(err)
	}
	defer listen.Close()

	for{
		fmt.Println("正在等待客户端连接")
		conn,err:=listen.Accept()
		if err!=nil{
			log.Panic(err)
		}
		fmt.Println("已有客户端连接",conn.RemoteAddr())
		//读取对方传来的数据
		request,err:=ioutil.ReadAll(conn)
		if err!=nil{
			log.Panic(err)
		}
		fmt.Printf("接收到的数据是：%s\n" ,request)
	}
}
