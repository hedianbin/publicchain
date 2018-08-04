package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

//这里已经有带有创世区块的区块链
type CLI struct{}

//输出说明书
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreatewallet -- 创建钱包")
	fmt.Println("\tcreateblockchain -address -- 创建创世区块")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tprintchain - 输出区块信息")
	fmt.Println("\tgetbalance -address Data -- 查询余额")
}

//验证命令行是否后面输入参数
func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

//命令行解析
func (cli *CLI) Run() {
	isValidArgs()
	//自定义cli命令
	createWlletCmd:=flag.NewFlagSet("createwallet",flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转帐源地址")
	flagTo := sendBlockCmd.String("to", "", "转帐目标地址")
	flagAmount := sendBlockCmd.String("amount", "", "转帐金额")
	flagGetBalanceData:=getBalanceCmd.String("address","","要查询余额的帐户")

	flagCreateBlockchainWithAddress := createBlockChainCmd.String("address", "", "创世区块的地址")
	//拿到第2个参数做判断

	switch os.Args[1] {
	//如果第2个参数输入的是addBlock就执行case "addBlock":里的代码
	case "send":
		//解析命令行并取出addBlock -data "liyuechun"第2个数组往后的所有数据。也就是取出-data后面的"liyuechun"
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		//解析命令行并取出printchain的所有数据"
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		//解析命令行并取出printchain的所有数据"
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		//解析命令行并取出printchain的所有数据"
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		//解析命令行并取出printchain的所有数据"
		err := createWlletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
	//判断addBlock参数后面的输入是否合法
	if sendBlockCmd.Parsed() { //如果解析成功
		//如果addBlock -data输入的是""空字符串
		if *flagFrom=="" || *flagTo=="" || *flagAmount == "" {
			fmt.Println("from,to,amount不能为空")
			//输出帮助信息
			printUsage()
			//退出程序
			os.Exit(1)
		}
		//将json转成数组
		from :=JSONToArray(*flagFrom)
		to :=JSONToArray(*flagTo)
		amount :=JSONToArray(*flagAmount)
		//遍历一下传过来的多笔交易中的from,to，看看地址是否有效
		for i:=0;i<len(from);i++{
			if !IsValidAddress([]byte(from[i])) || !IsValidAddress([]byte(to[i])){
				fmt.Println("地址无效，无法转帐...")
				printUsage()
				os.Exit(1)
			}
		}
		//执行send方法，将三个参数的值传进来
		cli.Send(from,to,amount)

	}
	//判断printchain参数是否解析
	if printChainCmd.Parsed() {
		cli.printchain()
	}
	//判断createblockchain参数是否解析
	if createBlockChainCmd.Parsed() {
		if !IsValidAddress([]byte(*flagCreateBlockchainWithAddress)) {
			//输出帮助信息
			fmt.Println("地址有误，无法创建创世区块")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}
	if getBalanceCmd.Parsed(){
		if !IsValidAddress([]byte(*flagGetBalanceData)){
			fmt.Println("地址有误，不能查询余额")
			printUsage()
			os.Exit(1)
		}
		cli.GetBalance(*flagGetBalanceData)
	}
	if createWlletCmd.Parsed(){
		//创建钱包
		cli.CreateWallet()
	}
}



