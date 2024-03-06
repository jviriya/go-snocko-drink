package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Hash(input, key string) string {
	keyForSign := []byte(key)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func DesEncryption(key, iv, plainText []byte) ([]byte, error) {

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData := PKCS5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}

func DesDecryption(key, iv, cipherText []byte) ([]byte, error) {

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func OpensslDecrypt(keyStr string, ivStr string, text string) string {

	key, _ := hex.DecodeString(keyStr)
	iv, _ := hex.DecodeString(ivStr)
	ciphertext, _ := hex.DecodeString(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext = PKCS7UnPadding(plaintext)
	return fmt.Sprintf("%s\n", plaintext)
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateKeyMD5WithRandom8Hex(params, secretKey string) string {
	b := make([]byte, 8/2)
	rand.Read(b)
	randomHex := fmt.Sprintf("%x", b)
	keyMD5 := GenMD5(fmt.Sprintf("%s%s", secretKey, params))

	return fmt.Sprintf("%s%s%s", keyMD5[:16], randomHex, keyMD5[16:])
}
