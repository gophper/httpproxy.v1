package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"sort"
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

	s1, _ := ReadFile("./data.txt")
	a, _ := EncryptAES([]byte(s1))
	fmt.Println(a)

	b, _ := DecryptAES(a)
	fmt.Println(string(b))
}

func TestFunc3(t *testing.T) {

	s := []byte{21, 3, 3, 0, 26, 154, 221, 161, 82, 167, 20, 30, 164, 149, 204, 2, 149, 170, 238, 236, 196, 41, 11, 168, 13, 102, 6, 72, 181, 81, 156}
	fmt.Println(string(s))
}

func TestFunc4(t *testing.T) {

	a := []int{16, 75, 227, 29, 174, 106, 199, 147, 7, 201, 195, 180, 163, 142, 128, 51, 61, 178, 149, 211, 185, 68, 160, 84, 167, 45, 251, 99, 89, 115, 30, 170, 168, 91, 1, 17, 130, 159, 215, 25, 50, 119, 82, 97, 222, 237, 155, 200, 8, 194, 87, 31, 44, 161, 60, 183, 120, 109, 59, 85, 86, 80, 204, 63, 18, 190, 137, 33, 132, 189, 187, 118, 111, 15, 151, 35, 146, 22, 141, 49, 182, 133, 153, 81, 231, 140, 93, 131, 176, 158, 172, 241, 69, 156, 165, 150, 37, 74, 240, 218, 19, 76, 0, 228, 217, 2, 144, 113, 139, 103, 107, 116, 6, 219, 9, 117, 197, 205, 88, 3, 246, 126, 244, 48, 78, 110, 58, 207, 235, 245, 214, 181, 206, 56, 114, 254, 143, 10, 166, 53, 95, 28, 198, 148, 127, 125, 247, 122, 42, 96, 253, 108, 248, 21, 123, 234, 83, 252, 64, 250, 162, 124, 223, 145, 72, 171, 192, 154, 230, 38, 41, 92, 209, 43, 225, 39, 179, 129, 46, 213, 164, 243, 5, 136, 232, 188, 40, 208, 239, 62, 13, 20, 233, 70, 212, 152, 52, 98, 157, 226, 4, 105, 196, 216, 101, 121, 191, 71, 236, 249, 102, 14, 242, 220, 193, 47, 23, 90, 255, 27, 55, 238, 66, 65, 104, 184, 135, 186, 73, 94, 210, 173, 175, 12, 32, 24, 229, 79, 224, 11, 221, 203, 54, 138, 67, 134, 34, 169, 57, 77, 26, 202, 112, 36, 100, 177}
	b := []int{102, 34, 105, 119, 200, 182, 112, 8, 48, 114, 137, 239, 233, 190, 211, 73, 0, 35, 64, 100, 191, 153, 77, 216, 235, 39, 250, 219, 141, 3, 30, 51, 234, 67, 246, 75, 253, 96, 169, 175, 186, 170, 148, 173, 52, 25, 178, 215, 123, 79, 40, 15, 196, 139, 242, 220, 133, 248, 126, 58, 54, 16, 189, 63, 158, 223, 222, 244, 21, 92, 193, 207, 164, 228, 97, 1, 101, 249, 124, 237, 61, 83, 42, 156, 23, 59, 60, 50, 118, 28, 217, 33, 171, 86, 229, 140, 149, 43, 197, 27, 254, 204, 210, 109, 224, 201, 5, 110, 151, 57, 125, 72, 252, 107, 134, 29, 111, 115, 71, 41, 56, 205, 147, 154, 161, 145, 121, 144, 14, 177, 36, 87, 68, 81, 245, 226, 183, 66, 243, 108, 85, 78, 13, 136, 106, 163, 76, 7, 143, 18, 95, 74, 195, 82, 167, 46, 93, 198, 89, 37, 22, 53, 160, 12, 180, 94, 138, 24, 32, 247, 31, 165, 90, 231, 4, 232, 88, 255, 17, 176, 11, 131, 80, 55, 225, 20, 227, 70, 185, 69, 65, 206, 166, 214, 49, 10, 202, 116, 142, 6, 47, 9, 251, 241, 62, 117, 132, 127, 187, 172, 230, 19, 194, 179, 130, 38, 203, 104, 99, 113, 213, 240, 44, 162, 238, 174, 199, 2, 103, 236, 168, 84, 184, 192, 155, 128, 208, 45, 221, 188, 98, 91, 212, 181, 122, 129, 120, 146, 152, 209, 159, 26, 157, 150, 135, 218}

	sort.Ints(a)
	sort.Ints(b)
	fmt.Println(a)
	fmt.Println(b)
}

func TestGenKey(t *testing.T) {

	//key256 := ""

	rand.Seed(time.Now().UnixNano())

	ks := make([]byte, 256)

	for index, _ := range ks {
		ks[index] = byte(index)
	}

	fmt.Println(ks)

	for key, _ := range ks {
		randIndex := rand.Intn(256)
		ks[key], ks[randIndex] = ks[randIndex], ks[key]
	}

	fmt.Println(ks)
	fmt.Println(string(ks))
	dst := make([]byte, 0)
	base64.RawStdEncoding.Encode(dst, ks)
	fmt.Println(dst)
	fmt.Println(string(dst))
	ded := make([]byte, 1000)
	base64.RawStdEncoding.Decode(ded, dst)
	fmt.Println(ded)
}
