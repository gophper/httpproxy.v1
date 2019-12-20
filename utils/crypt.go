package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"httpproxy.v1/config"
	"io/ioutil"
	"os"
)

var ErrEmptyKey = errors.New("short write")
var keyStock = make(map[byte]byte, 256)
var valStock = make(map[byte]byte, 256)

func initCrypt() error {
	encodeString := config.GetConfig("sys", "key")
	if encodeString == "" {
		return ErrEmptyKey
	}
	kss, err := base64.StdEncoding.DecodeString(string(encodeString))
	if err != nil {
		return err
	}
	for key, value := range kss {
		keyStock[byte(key)] = value
	}
	for key, value := range keyStock {
		valStock[value] = byte(key)
	}
	return nil
}

func EncryptAES(src []byte) ([]byte, error) {
	var ok bool
	for key, value := range src {
		src[key], ok = keyStock[value]
		if !ok {
			fmt.Println(keyStock, value, key)
			panic("encode key not found")
		}
	}
	return src, nil
}

func DecryptAES(src []byte) (dst []byte, err error) {
	var ok bool
	for key, value := range src {
		src[key], ok = valStock[value]
		if !ok {
			fmt.Println(valStock, value, key)
			panic("decode key not found")
		}
	}
	return src, nil
}

func ReadFile(filename string) (content []byte, err error) {
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
