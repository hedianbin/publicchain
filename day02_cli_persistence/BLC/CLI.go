package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)
//这里已经有带有创世区块的区块链
type CLI struct {}
//输出说明书
func printUsage()  {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address -- 创建创世区块")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain - 输出区块信息")
}
//验证命令行是否后面输入参数
func isValidArgs()  {
	if len(os.Args)<2{
		printUsage()
		os.Exit(1)
	}
}

//参数是addblock后添加新区块到数据库中
func (cli *CLI) addBlock(txs []*Transaction)  {
	if DBExists()==false{
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain:=GetBlockchainObject()
	defer blockchain.DB.Close()
	blockchain.AddBlockToBlockchain(txs)
}
//打印输出数据库中所有区块信息
func (cli *CLI) printchain()  {
	if DBExists()==false{
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain:=GetBlockchainObject()
	defer blockchain.DB.Close()
	blockchain.Printchain()
}
//命令行解析
func (cli *CLI) Run()  {
	isValidArgs()
	//自定义cli命令
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd:=flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockChainCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)
	flagAddBlockData := addBlockCmd.String("data", "http://www.baidu.com", "交易数据....")
	flagCreateBlockchainWithAddress:=createBlockChainCmd.String("address","","创世区块的地址")
	//拿到第2个参数做判断

	switch os.Args[1] {
	//如果第2个参数输入的是addBlock就执行case "addBlock":里的代码
	case "addblock":
		//解析命令行并取出addBlock -data "liyuechun"第2个数组往后的所有数据。也就是取出-data后面的"liyuechun"
		err:=addBlockCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "printchain":
		//解析命令行并取出printchain的所有数据"
		err:=printChainCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "createblockchain":
		//解析命令行并取出printchain的所有数据"
		err:=createBlockChainCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	//判断addBlock参数后面的输入是否合法
	if addBlockCmd.Parsed(){ //如果解析成功
		//如果addBlock -data输入的是""空字符串
		if *flagAddBlockData == ""{
			//输出帮助信息
			printUsage()
			//退出程序
			os.Exit(1)
		}
		//如果获取到-data后面的内容，就添加新区块
		cli.addBlock([]*Transaction{})
	}
	//判断printchain参数是否解析
	if printChainCmd.Parsed(){
		cli.printchain()
	}
	//判断createblockchain参数是否解析
	if createBlockChainCmd.Parsed(){
		if *flagCreateBlockchainWithAddress==""{
			//输出帮助信息
			fmt.Println("地址不能为空.......")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}
}
//创建创世区块
func (cli *CLI) createGenesisBlockchain(address string) {
	CreateBlockchainWithGenesisBlock(address)
}