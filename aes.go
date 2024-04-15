package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
)

// 解密
func decryptFunc(cipherText, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	// 去除填充
	padLength := int(cipherText[len(cipherText)-1])
	return string(cipherText[:len(cipherText)-padLength]), nil
}

// 加密
func encryptFunc(data interface{}, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 如果数据是字符串类型，将其转换为字节数组
	var plaintext []byte
	switch v := data.(type) {
	case string:
		plaintext = []byte(v)
	default:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		plaintext = jsonData
	}

	// 填充明文数据
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// 创建CBC模式的加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密数据
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// PKCS7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// 初始化向量
func intAesIV() string {
	var a string
	for i := 0; i < 16; i++ {
		a += string(rune(i + 112))
	}
	return a
}
