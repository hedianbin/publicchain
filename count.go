package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"log"
	"bufio"
	"io"
)

func main() {
	filepaths := []string{
		"D:\\gopath\\src\\dotcoin\\addr",
		"D:\\gopath\\src\\dotcoin\\base58",
		"D:\\gopath\\src\\dotcoin\\chain",
		"D:\\gopath\\src\\dotcoin\\cli",
		"D:\\gopath\\src\\dotcoin\\config",
		"D:\\gopath\\src\\dotcoin\\connx",
		"D:\\gopath\\src\\dotcoin\\const",
		"D:\\gopath\\src\\dotcoin\\logx",
		"D:\\gopath\\src\\dotcoin\\mempool",
		"D:\\gopath\\src\\dotcoin\\merkle",
		"D:\\gopath\\src\\dotcoin\\mining",
		"D:\\gopath\\src\\dotcoin\\peer",
		"D:\\gopath\\src\\dotcoin\\protocol",
		"D:\\gopath\\src\\dotcoin\\server",
		"D:\\gopath\\src\\dotcoin\\storage",
		"D:\\gopath\\src\\dotcoin\\sync",
		"D:\\gopath\\src\\dotcoin\\util",
		"D:\\gopath\\src\\dotcoin\\wallet",
	}
	var total int
	for _,filepath:=range filepaths{
		num, _ := GetAllFile(filepath)
		total+=num
	}

	fmt.Println(total)
}

func GetAllFile(pathname string) (int, error) {
	rd, err := ioutil.ReadDir(pathname)
	var total int
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			GetAllFile(pathname + fi.Name() + "\\")
		} else {
			fmt.Println(fi.Name())
			num := ReadFile(pathname + "\\" + fi.Name())
			total += num
		}
	}
	return total, err
}
func ReadFile(filename string) int {
	var num int
	fi, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		_, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		num++
	}
	return num
}
