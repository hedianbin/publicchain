package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

//将int64转换为字节数组
func IntToHex(num int64) []byte  {
	//新建缓冲区用来存转换的数据
	buff:=new(bytes.Buffer)
	//将num转换为字节数组存到buff中
	err:=binary.Write(buff,binary.BigEndian,num)
	if err!=nil{
		log.Panic(err)
	}
	//返回转换后的字节数组
	return buff.Bytes()
}
