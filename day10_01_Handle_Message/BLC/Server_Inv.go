package BLC

type Inv struct {
	AddrFrom string //当前节点自己的地址
	Type string //要发送的数据总数
	Items [][]byte  //要传递的数据的hash
}

