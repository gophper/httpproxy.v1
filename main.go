package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

func main() {

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
	dst := make([]byte,0)
	base64.RawStdEncoding.Encode(dst,ks)
	fmt.Println(dst)
	fmt.Println(string(dst))
	ded := make([]byte,1000)
	base64.RawStdEncoding.Decode(ded,dst)
	fmt.Println(ded)
}


