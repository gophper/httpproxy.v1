package crypt
import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"strings"
)

func Md5Encode(v string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(v))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func Hmac(key, data string) string {
	hmacHash := hmac.New(md5.New, []byte(key))
	hmacHash.Write([]byte(data))
	return hex.EncodeToString(hmacHash.Sum([]byte("")))
}

// aes encrpty alg with mode of CBC
func AesCBCEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// aes decrpty alg with mode of CBC
func AesCBCDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// padding deal
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// unpadding deal
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// gzflate decode
func Gzdecode(data string) string {
	if data == "" {
		return ""
	}
	r := flate.NewReader(strings.NewReader(data))
	defer r.Close()
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}
	return string(out)
}

// gzflate encode
func Gzencode(data string, level int) []byte {
	if data == "" {
		return []byte{}
	}
	var bufs bytes.Buffer
	w, _ := flate.NewWriter(&bufs, level)
	w.Write([]byte(data))
	w.Flush()
	w.Close()
	return bufs.Bytes()
}

func AesEncode(data string, aesKey []byte) string {
	if data == "" {
		return data
	}

	gzEncodeData := Gzencode(data, 9)
	aesEncryptData, err := AesCBCEncrypt(gzEncodeData, aesKey)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(aesEncryptData)
}

func AesDecode(data string, aesKey []byte) string {
	if data == "" {
		return data
	}

	crypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		println(err.Error())
		return ""
	}

	aesDecodeData, err := AesCBCDecrypt(crypted, aesKey)
	if err != nil {
		println(err.Error())
		return ""
	}

	return Gzdecode(string(aesDecodeData))
}
