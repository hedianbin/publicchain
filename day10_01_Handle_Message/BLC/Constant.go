package BLC

//数据库名字
const dbName = "blockchain_%s.db"

//表的名字
const blockTableName = "blocks"

//windows电脑=0,linux和macos=1
const SYSTEM_SELECT = 0

//命令的长度
const COMMAND_LENGTH = 14

//定义程序版本
const NODE_VERSION = 1

//具体的命令version
const COMMAND_VERSION = "version"

//拿到所有的hash
const COMMAND_GETBLOCKHASHS = "getblockhashs"

//类型
const BLOCK_TYPE = "block"
const TX_TYPE = "tx"

//发送所有区块hash
const COMMAND_INV = "inv"

//发送所要的区块数据
const COMMAND_BLOCKDATA = "blockdata"

//发送需要的区块hash给对方，向对方索要区块
const COMMAND_GETBLOCKDATA = "getblockdata"