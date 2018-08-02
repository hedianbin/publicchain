package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

func main() {
	// 在当前目录下打开my.db数据库
	// 如果不存在就创建它
	db,err:=bolt.Open("my.db",0600,nil)
	if err!=nil{
		log.Fatal(err)
	}
	defer db.Close()
	err=db.View(func(tx *bolt.Tx) error {
		//打开表
		b:=tx.Bucket([]byte("BlockBucket"))
		//如果b不为空，就说明表创建成功。可以往里存储数据
		if b!=nil{
			//获取key为"l"对应的value值
			data:=b.Get([]byte("l"))
			fmt.Printf("%s\n",data)
			data=b.Get([]byte("ll"))
			fmt.Printf("%s\n",data)
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
}
