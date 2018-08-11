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
		//管对方要所有区块hash,参数传对方地址
		sendGetBlockHashs(version.AddrFrom)
	}
}

//找到对方要的所有区块hash,发送给对方
func handleGetBlockHashs(request []byte,bc *Blockchain)  {
	//取命令
	command:=bytesToCommand(request[:COMMAND_LENGTH])
	getblocksBytes:=request[COMMAND_LENGTH:] //取数据

	//将数据进行反序列化
	var getblocks GetBlocks
	reader:=bytes.NewReader(getblocksBytes)
	decoder:=gob.NewDecoder(reader)
	err:=decoder.Decode(&getblocks)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("接收到了来自%s的命令是：%s\n",getblocks.AddrFrom,command)
	//查询自己的数据表，将所有区块的hash拼接，发送给对方
	blockHashes:=bc.getBlocksHashes()
	//发送所有的区块hash给对方
	sendInv(getblocks.AddrFrom,BLOCK_TYPE,blockHashes)
}

//处理对方发来的所有区块hash
func handleInv(request []byte,bc *Blockchain)  {
	command:=bytesToCommand(request[:COMMAND_LENGTH])
	invBytes:=request[COMMAND_LENGTH:]

	var inv Inv
	//反序列化
	reader:=bytes.NewReader(invBytes)
	decoder:=gob.NewDecoder(reader)
	err:=decoder.Decode(&inv)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("接收到了来自：%s,传来的%s\n",inv.AddrFrom,command)
	if inv.Type==BLOCK_TYPE{
		//发送获取对应的数据请求
		hash:=inv.Items[0]//第一个创建区块的hash,block0
		//给对方发送过去我要的hash,然后对方会给我我要的区块
		sendGetBlockData(inv.AddrFrom,BLOCK_TYPE,hash)
		//判断，除了block0的hash以外，把其它hash都添加到blocksArray数组中
		if len(inv.Items)>=1{
			blocksArray=inv.Items[1:]
		}
	}else if inv.Type==TX_TYPE{

	}
}

//根据请求hash拿区块数据
func handleGetBlockData(request []byte,bc *Blockchain)  {
	command:=bytesToCommand(request[:COMMAND_LENGTH])
	getDataBytes:=request[COMMAND_LENGTH:]

	//反序列化
	var getData GetData
	reader:=bytes.NewReader(getDataBytes)
	decoder:=gob.NewDecoder(reader)
	err:=decoder.Decode(&getData)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("接收到了来自：%s,传来的命令：%s\n",getData.AddrFrom,command)
	if getData.Type==BLOCK_TYPE{
		//根据hash找到对应区块
		block:=bc.GetBlockByHash(getData.Hash)
		sendBlock(getData.AddrFrom,block)
	}else if getData.Type==TX_TYPE{

	}
}
//拿到了想要的区块数据，存储到数据库中
func handleBlockData(request []byte,bc *Blockchain)  {
	command:=bytesToCommand(request[:COMMAND_LENGTH])
	blockDataBytes:=request[COMMAND_LENGTH:]
	//反序列化
	var blockData BlockData
	reader:=bytes.NewReader(blockDataBytes)
	decoder:=gob.NewDecoder(reader)
	err:=decoder.Decode(&blockData)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("收到来自：%s传来的命令:%s\n",blockData.AddrFrom,command)

	//取出block数据，存入自己的数据库中
	blockBytes:=blockData.Block
	block:=DeserializeBlock(blockBytes)
	//存入到本地的数据库中
	bc.AddBlock(block)

	//如果要添加的区块没了，就更新utxoset表
	if len(blocksArray)==0{
		utxoset:=UTXOSet{bc}
		utxoset.ResetUTXOSet()
		fmt.Println("所有区块已经同步完成")
	}
	if len(blocksArray)>0{
		fmt.Printf("正在更新Height为%d的区块。。。。\n",block.Height)
		//发送请求，继续获取下一个区块数据
		sendGetBlockData(blockData.AddrFrom,BLOCK_TYPE,blocksArray[0])
		//删除刚才获取的那个区块的hash
		blocksArray=blocksArray[1:]
	}

}