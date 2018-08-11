package BLC

import (
	"crypto/sha256"
	"math"
)

//第一步，创建结构体对象，表示节点和树
//这是节点
type MerkleNode struct {
	LeftNode *MerkleNode //左子树
	RightNode *MerkleNode //右子树
	DataHash []byte //数据hash
}

//树,只有根节点
type MerkleTree struct {
	RootNode *MerkleNode
}

//第二步，给一个左右节点，生成一个新的节点
func NewMerkleNode(leftNode,rightNode *MerkleNode,txHash []byte) *MerkleNode {
	//1.创建爹，现在爹还没儿子
	mNode:=&MerkleNode{}
	//2.如果爹没有左儿子。没有右儿子，那么爹自己就是儿子。
	if leftNode==nil && rightNode==nil{
		//爹自己就是儿子，生成自己的哈希值
		hash:=sha256.Sum256(txHash)
		mNode.DataHash=hash[:]
	}else{
		//否则如果有左儿子和右儿子，但他俩现在还没爹。
		prevHash:=append(leftNode.DataHash,rightNode.DataHash...)
		hash:=sha256.Sum256(prevHash) //他俩加起来再一起哈希
		mNode.DataHash=hash[:] //帮爹生成哈希值
	}
	mNode.LeftNode=leftNode //左儿子认爹
	mNode.RightNode=rightNode //右儿子认爹
	//返回有左右儿子的爹，爹也有自己的哈希
	return mNode
}

//生成MerkleTree
func NewMerkleTree(txHashData [][]byte) *MerkleTree  {
	//1.创建一个数组，用于存储node节点
	var nodes []*MerkleNode
	//2.判断交易量的奇偶性
	if len(txHashData)%2!=0{
		//奇数，复制最后一个
		txHashData=append(txHashData,txHashData[len(txHashData)-1])
	}
	//3.创建一排叶子节点，遍历交易hash数据
	for _,datum:=range txHashData{
		//叶子节点左右节点都为nil,再传进去交易tx序列化数据
		node:=NewMerkleNode(nil,nil,datum)
		//将生成好的叶子节点都追加到nodes数组中
		nodes=append(nodes,node)
	}
	//4.生成树其他的节点
	count:=GetCircleCount(len(nodes))
	//生成其它节点，直到生成到最后的根节点
	for i:=0;i<count ;i++  {
		var newLevel []*MerkleNode
		//两两哈希
		for j:=0;j<len(nodes);j+=2{
			//两两哈希，生成爹，爹就有左右儿子了
			node:=NewMerkleNode(nodes[j],nodes[j+1],nil)
			//将生成的新节点追加到newLevel数组中
			newLevel=append(newLevel,node)
		}
		//先判断newLevel的奇偶性
		if len(newLevel)%2!=0{ //如果是奇数
			newLevel=append(newLevel,newLevel[len(newLevel)-1]) //最后一个生成副本
		}
		nodes=newLevel
	}
	//拿到rootNode
	mTree:=&MerkleTree{nodes[0]}
	//返回根节点
	return mTree
}

//统计几层
func GetCircleCount(len int) int {
	count:=0
	for{
		//计算2的几次方>=len
		if int(math.Pow(2,float64(count)))>=len{
			return count //如果找到了层数就返回
		}
		//没找到层数，就继续++
		count++
	}
}