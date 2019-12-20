package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestBase64(t *testing.T) {
	input := []byte("hello world")
	encodeString := base64.StdEncoding.EncodeToString(input)
	fmt.Println(encodeString)
	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(decodeBytes))

	// 如果要用在url中，需要使用URLEncoding
	uEnc := base64.URLEncoding.EncodeToString([]byte(input))
	fmt.Println(uEnc)

	uDec, err := base64.URLEncoding.DecodeString(uEnc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(uDec))
}

//func TestAES(t *testing.T) {
//	d := []byte("hello,ase")
//	key := []byte("hgfedcba87654321")
//	fmt.Println("加密前:", string(d))
//	x1, err := EncryptAES(d, key)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println("加后密:", string(x1))
//	x2, err := DecryptAES(x1, key)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println("解密后:", string(x2))
//}

func TestJm(t *testing.T) {
	s := make([]int, 256)
	for key, _ := range s {
		s[key] = key
	}
	fmt.Println(s)
	rand.Seed(time.Now().UnixNano())

	for key, _ := range s {
		i := s[rand.Intn(256)]
		s[key], s[i] = s[i], s[key]
	}

	bs := make([]byte, 256)
	for key, value := range s {
		bs[key] = byte(value)
	}
	fmt.Println(bs)
	encodeString := base64.StdEncoding.EncodeToString(bs)
	fmt.Println(encodeString)

	WriteFile("./test.txt", encodeString)
}

func TestFunc2(t *testing.T) {

	s1, _ := ReadFile_v1("./data.txt")
	a, _ := EncryptAES([]byte(s1))
	fmt.Println(a)

	b, _ := DecryptAES(a)
	fmt.Println(string(b))
}

func TestFunc3(t *testing.T)  {

	s := []byte{21,3,3,0,26,154,221,161,82,167,20,30,164,149,204,2,149,170,238,236,196,41,11,168,13,102,6,72,181,81,156}
	fmt.Println(string(s))
}
