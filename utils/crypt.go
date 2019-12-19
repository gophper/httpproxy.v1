package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var keyStock []byte
var valStock = make([]byte,256)

func init() {
	encodeString, err := ReadFile_v1("E:\\code\\go\\src\\gostudy\\httpproxy\\utils\\test.txt")

	// 对上面的编码结果进行base64解码
	keyStock, err = base64.StdEncoding.DecodeString(string(encodeString))
	if err != nil {
		log.Fatalln(err)
	}
	for key, value := range keyStock {
		valStock[value] = byte(key)
	}
}

func EncryptAES(src []byte) ([]byte, error) {
	//return src, nil
	for key, value := range src {
		src[key] = keyStock[value]
	}
	return src, nil
}

func DecryptAES(src []byte) (dst []byte, err error) {
	//return src, nil
	for key, value := range src {
		src[key] = valStock[value]
	}
	return src, nil
}

func ReadFile_v1(filename string) (content []byte, err error) {
	fileObj, err := os.Open(filename)
	if err != nil {
		fmt.Println("os open error:", err)
		return
	}
	defer fileObj.Close()
	content, err = ioutil.ReadAll(fileObj)
	if err != nil {
		fmt.Println("ioutil.ReadAll error:", err)
		return
	}
	return
}

func WriteFile(filename, data string) {
	var (
		err error
	)
	fileObj, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = fmt.Fprintf(fileObj, data)
}
