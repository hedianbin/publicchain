package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"strings"
	"encoding/gob"
	"fmt"
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

//标准的JSON字符串转数组,传一个json string字符串
func JSONToArray(jsonString string) []string {
	var repStr string
	if SYSTEM_SELECT==0{
		repStr=strings.Replace(jsonString,"'[","["+"\"",-1)
		repStr=strings.Replace(repStr,"]'","\""+"]",-1)
		repStr=strings.Replace(repStr,",","\""+","+"\"",-1)
	}
	//定义一个数组，用来装转换过来的数组
	var sArr []string
	//将json字符串转换为数组存储到sArr中。
	if err:=json.Unmarshal([]byte(repStr),&sArr);err!=nil{
		log.Panic(err)
	}
	return sArr
}

//字节数组反转
func ReverseBytes(data []byte) {
	for i,j:= 0,len(data)-1; i <j ;i,j=i+1,j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

//将对象进行序列化
func gobEncode(data interface{}) []byte  {
	var buff bytes.Buffer
	encoder:=gob.NewEncoder(&buff)
	err:=encoder.Encode(data)
	if err!=nil{
		log.Panic(err)
	}
	return buff.Bytes()
}

func commandToBytes(command string) []byte {
	var bytes [COMMAND_LENGTH]byte //定义一个长度为12的数组
	//遍历每个命令字符
	for i,c:=range command {
		bytes[i]=byte(c) //将命令字符转成byte赋值给数组，[v,e,r,s,i,o,n,0,0,0,0,0]
	}
	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte //用来存储拿到的命令
	for _,b:=range bytes{ //遍历命令的12个字节，把0删除，只保留前面的命令
		if b!=0x00{
			command=append(command,b)
		}
	}
	//返回拿到的命令
	return fmt.Sprintf("%s",command)
}