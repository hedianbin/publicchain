package BLC


import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)

//版本号
const version = byte(0x00)
//地址的长度
const addressCheckSumLen =4

//定义一个钱包结构
type Wallet struct {
	//1.私钥
	PrivateKey ecdsa.PrivateKey
	//2.公钥
	PublicKey []byte
}

//创建钱包对象
func NewWallet() *Wallet {
	//生成私钥和公钥
	privateKey,publicKey:=newKeyPair()
	//创建钱包对象，得到了钱包，里面有生成的私钥和公钥
	wallet:=&Wallet{privateKey,publicKey}
	//返回创建的钱包
	return wallet
}

//2.产生一对密钥，私钥和公钥
func newKeyPair() (ecdsa.PrivateKey,[]byte) {
	/*
	1.根据椭圆曲线加密产生随机私钥，
	2.根据私钥产生公钥
	椭圆：ellipse
	曲线：curve
	椭圆曲线加密：(ECC:ellipse curve Cryptography),非对称加密

	*/
	//根据椭圆曲线加密算法，得到一个椭圆曲线值
	curve:=elliptic.P256()
	//生成私钥,rand.Reader为随机数
	privateKey,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{
		log.Panic(err)
	}
	//用产生的私钥的X轴的值和Y轴的值产生公钥
	publicKey:=append(privateKey.PublicKey.X.Bytes(),privateKey.PublicKey.Y.Bytes()...)
	//返回私钥和公钥
	return *privateKey,publicKey
}


//生成公钥哈希
func PubKeyHash(publicKey []byte) []byte  {
	//1.将原始公钥经过一次sha256,再一次160,得到公钥哈希
	hasher:=sha256.New() //先创建一个sha256的hash对象
	hasher.Write(publicKey)//然后将公钥写到hasher中
	hash1:=hasher.Sum(nil) //然后再将写入的hasher生成sha256的hash

	//2.再一次160,原来和sha256一样
	hasher2:=ripemd160.New()
	hasher2.Write(hash1)
	hash2:=hasher2.Sum(nil)
	//3.返回160
	return hash2
}

//生成较验码，并取前4位返回
func CheckSum(payload []byte) []byte {
	//第一次sha256
	firstHash:=sha256.Sum256(payload)
	//第二次sha256
	secondHash:=sha256.Sum256(firstHash[:])
	//取前4位
	return secondHash[:addressCheckSumLen]
}

//根据公钥生成地址,属于Wallet的方法，就能直接用公钥和私钥
func (w *Wallet) GetAddress() []byte {
	//1.将原始公钥经过一次sha256,再一次160,得到公钥哈希160的
	pubKeyHash:=PubKeyHash(w.PublicKey)
	address:=GetAddressByPubKeyHash(pubKeyHash)
	return address
}
//计算钱包地址
func GetAddressByPubKeyHash(pubKeyHash []byte) []byte {
	//2.拼接版本号+公钥哈希
	versioned_payload:=append([]byte{version},pubKeyHash...)
	//拿拼接后的字节数组，经过两次sha256,然后取前4位得到较验码
	checkSumBytes:=CheckSum(versioned_payload)
	//3.将版本号+公钥哈希+较验码拼接起来，得到拼后的字节数组
	full_payload:=append(versioned_payload,checkSumBytes...)
	//再经过base58，拿到钱包地址
	address:=Base58Encode(full_payload)
	return address
}
//校验地址是否有效
func IsValidAddress(address []byte) bool {
	//1.Base58解码
	full_payload:=Base58Decode(address)
	//2.获取地址当中携带的4位较验码
	checkSumBytes:=full_payload[len(full_payload)-addressCheckSumLen:]
	version:=full_payload[0] //取得解密后的base58第1位字节为0，[0 67 217 197 44 207 52]
	//获取版本号+公钥哈希
	versioned_payload:=full_payload[1:len(full_payload)-addressCheckSumLen]
	//3.重新用版本号+公钥哈希生成一次较验码
	checkSumBytes2:=CheckSum(append([]byte{version},versioned_payload...))
	//4.比较携带的验证码和生成的一样就没问题
	return bytes.Compare(checkSumBytes,checkSumBytes2)==0
}