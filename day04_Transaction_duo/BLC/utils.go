package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"strings"
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
//标准的JSON字符串转数组,传一个json string字符串
func JSONToArray(jsonString string) []string {
	repStr:=strings.Replace(jsonString,"'[","["+"\"",-1)
	repStr=strings.Replace(repStr,"]'","\""+"]",-1)
	repStr=strings.Replace(repStr,",","\""+","+"\"",-1)
	//定义一个数组，用来装转换过来的数组
	var sArr []string
	//将json字符串转换为数组存储到sArr中。
	if err:=json.Unmarshal([]byte(repStr),&sArr);err!=nil{
		log.Panic(err)
	}
	return sArr
}