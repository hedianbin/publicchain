package BLC

type BlockData struct {
	AddrFrom string //当前节点自己的地址
	Block []byte //要传递序列化后的区块数据
}

