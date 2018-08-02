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
	fmt.Println("\tcreateblockchain -address -- 创建创世区块")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tprintchain - 输出区块信息")
}

//验证命令行是否后面输入参数
func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

//打印输出数据库中所有区块信息
func (cli *CLI) printchain() {
	if DBExists() == false {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain := GetBlockchainObject()
	defer blockchain.DB.Close()
	blockchain.Printchain()
}

//命令行解析
func (cli *CLI) Run() {
	isValidArgs()
	//自定义cli命令
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转帐源地址")
	flagTo := sendBlockCmd.String("to", "", "转帐目标地址")
	flagAmount := sendBlockCmd.String("amount", "", "转帐金额")

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
	default:
		printUsage()
		os.Exit(1)
	}
	//判断addBlock参数后面的输入是否合法
	if sendBlockCmd.Parsed() { //如果解析成功
		//如果addBlock -data输入的是""空字符串
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			//输出帮助信息
			printUsage()
			//退出程序
			os.Exit(1)
		}
		//将json转成数组
		from :=JSONToArray(*flagFrom)
		to :=JSONToArray(*flagTo)
		amount :=JSONToArray(*flagAmount)
		//执行send方法，将三个参数的值传进来
		cli.Send(from,to,amount)
	}
	//判断printchain参数是否解析
	if printChainCmd.Parsed() {
		cli.printchain()
	}
	//判断createblockchain参数是否解析
	if createBlockChainCmd.Parsed() {
		if *flagCreateBlockchainWithAddress == "" {
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
	blockchain:=CreateBlockchainWithGenesisBlock(address)
	defer blockchain.DB.Close()
}

//转帐
func (cli *CLI) Send(from []string,to []string,amount []string)  {
	if !DBExists(){
		fmt.Println("数据库不存在......")
		os.Exit(1)
	}
	//拿到了带有最新区块hash和db对象的Blockchain对象
	bc:=GetBlockchainObject()
	//关闭数据库连接
	defer bc.DB.Close()
	//挖掘新区块，带Transaction交易的
	bc.MineNewBlock(from,to,amount)
}
