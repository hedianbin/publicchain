package main

import (
	"flag"
	"os"
	"log"
	"fmt"
)

func main() {
	isValidArgs()
	//自定义cli命令
	addBlockCmd:=flag.NewFlagSet("addBlock",flag.ExitOnError)
	flagAddBlockData:=addBlockCmd.String("data","http://www.baidu.com","交易数据.....")
	printChainCmd:=flag.NewFlagSet("printchain",flag.ExitOnError)
	switch os.Args[1] {
	case "addBlock":
		err:=addBlockCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "printchain":
		err:=printChainCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	//判断addBlock参数后面的输入是否合法
	if addBlockCmd.Parsed(){
		if *flagAddBlockData==""{
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*flagAddBlockData)
	}
	if printChainCmd.Parsed(){
		fmt.Println("输出区块所有数据.......")
	}
}

func printUsage()  {
	fmt.Println("Usage:")
	fmt.Println("\taddBlock -data DATA -- 交易数据")
	fmt.Println("\tprintchain - 输出区块信息")

}

func isValidArgs()  {
	if len(os.Args)<2{
		printUsage()
		os.Exit(1)
	}
}