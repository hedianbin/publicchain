package BLC

import (
	"math/big"
	"bytes"
)

//base64

/*
ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/
0(零)，O(大写的o)，I(大写的i)，l(小写的L),+,/
 */
//字母表格，最终会展示的字符
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

//字节数组转Base58，加密，input为输入的二进制
func Base58Encode(input []byte) []byte {
	var result []byte
	//x生成的是这样一串4864296617342959715134812811642943631993346213843746529716
	//并且每次生成的都不一样
	//将input转换二进制并按整数风格存储
	x := big.NewInt(0).SetBytes(input)
	//base生成的就是58,创建了一个大数，存储的就是我们的二进制数
	base := big.NewInt(int64(len(b58Alphabet)))
	//zero就是0
	zero := big.NewInt(0)
	//创建big.int对象，取地址
	mod := &big.Int{}
	/*
	如要将1234转换为58进制；
	第一步：1234除于58，商21，余数为16，查表得H
	第二步：21除于58，商0，余数为21，查表得N
	所以得到base58编码为：NH
	*/
	for x.Cmp(zero) != 0 {
		//x除base，取余数存到mod中,商赋值给x
		x.DivMod(x, base, mod)
		//拿着mod中的余数当作b58Alphabet的下标去查表得到字典里的对应的值，追加到result中。直到x也就是商除到0为止退出循环
		//拿到的result是倒过来的
		result = append(result, b58Alphabet[mod.Int64()])
	}
	//将拿到的result反转
	//因为之前先附加低位的，后附加高位的，所以需要翻转
	ReverseBytes(result)

	//因为如果高位为0，0除任何数为0，可以直接设置为‘1’
	for b := range input {
		if b ==0x00 {
			//如果遍历的input的value字节数组是0就在高位补1，并不断的追加
			result = append([]byte{b58Alphabet[0]}, result...)
		}else{
			break
		}

	}
	//返回编码后的Base58
	return result

}

//Base58转字节数组，解密
func Base58Decode(input []byte) []byte {
	//result是big.int值为0
	result := big.NewInt(0) //初始化为0
	zeroBytes := 0 //int型0,记数
	//得到高位0的位数
	for b := range input {
		if b == 0x00 { //循环
			zeroBytes++//记数叠加,判断一下多少位
		}
	}
	payload := input[zeroBytes:] //取出要解码的字节最后几位
	for _,b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b) //字母表格
		result.Mul(result, big.NewInt(58)) //乘法
		result.Add(result, big.NewInt(int64(charIndex))) //加法
	}

	decoded := result.Bytes() //解码
	//一部分没有0，我们要统一补0，叠加
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}


