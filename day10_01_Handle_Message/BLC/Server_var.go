package BLC

//全节点地址
var knowNodes=[]string{"localhost:3000"}

//当前节点自己的地址
var nodeAddress string

//记录应该同步但尚未同步的区块hash
var blocksArray [][]byte