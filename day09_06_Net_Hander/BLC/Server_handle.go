package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)

//处理接收到的version数据
func handleVersion(request []byte,bc *Blockchain)  {
	/*
	1.从request中获取版本的数据：[]byte
	2.反序列化---Version{}
	3.操作bc,获取自己的最后block的height
	4.跟对方比较
	*/

	//1.从request中获取版本的数据：[]byte
	versionBytes:=request[COMMAND_LENGTH:]
	//2.反序列化---Version{}
	var version Version
	reader:=bytes.NewReader(versionBytes)
	decoder:=gob.NewDecoder(reader)
	err:=decoder.Decode(&version)
	if err!=nil{
		log.Panic(err)
	}
	//3.操作bc,获取自己的最后block的height
	selfHeight:=bc.GetBestHeight() //自己的高度
	foreignerBestHeight:=version.BestHeight //对方发来的高度
	fmt.Printf("接收到%s传来的版本高度%d\n",version.AddrFrom,foreignerBestHeight)

	//4.跟对方比较
	if selfHeight>foreignerBestHeight{
		//发送版本对象给对方
		sendVersion(version.AddrFrom,bc)
	}else{
		fmt.Println("我的高度没有你高，给我看看你的数据。。。。")
	}
}
