package BLC

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
)

//节点服务端启动，可以接收其他节点传来的数据
func startServer(nodeID string, mineAddress string) {
	//当前节点自己的地址
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	//监听地址
	listener, err := net.Listen("tcp", nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	//延迟关闭连接
	defer listener.Close()

	//拿到bc对象
	bc := GetBlockchainObject(nodeID)
	//如果不是全节点，比如是钱包节点或矿工节点，就给全节点发送一个消息
	if nodeAddress != knowNodes[0] {
		//给全节点发送版本信息
		sendVersion(knowNodes[0], bc)
	}
	//如果是全节点。等待客户端连入,当然不是全节点也可以等待连入
	for {
		//等待客户端连入，阻塞了
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("发送方已经连入：", conn.RemoteAddr())
		//处理接收到的数据
		go handleConnection(conn, bc)
	}
}

//按照拿到的命令，处理不同的请求
func handleConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn) //读到传来的数据
	if err != nil {
		log.Panic(err)
	}
	//从包里拿到命令
	command := bytesToCommand(request[:COMMAND_LENGTH])
	fmt.Printf("接收到的命令是：%s\n", command)

	switch command {
	case COMMAND_VERSION:
		//此处是处理接收到版本数据
		handleVersion(request,bc)
	case COMMAND_GETBLOCKHASHS:
		//看我的所有区块hashs
		handleGetBlockHashs(request,bc)
	case COMMAND_INV:
		//处理对方发过来的所有区块hash
		handleInv(request,bc)
	case COMMAND_GETBLOCKDATA:
		//根据传过来的hash拿区块数据
		handleGetBlockData(request,bc)
	case COMMAND_BLOCKDATA:
		//已经发来了真正的区块，把区块存储到数据库中
		handleBlockData(request,bc)
	default:
		fmt.Println("读不懂的命令。。。")
	}
	//关闭连接
	defer conn.Close()
}


