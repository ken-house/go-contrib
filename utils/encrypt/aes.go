package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// AesEncrypt AES加密
func AesEncrypt(contentStr string, key string, iv string) (string, error) {
	// 1.创建AES加密对象
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", nil
	}
	// 2.CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))

	// 3.对原文分块并设置PKCS5Padding填充模式对数据进行填充
	// a.获取数据块大小为8位
	blockSize := block.BlockSize()
	// b.获取要填充的位数
	padNum := blockSize - len(contentStr)%blockSize
	// c.生成要填充的字节数组
	padByte := bytes.Repeat([]byte{byte(padNum)}, padNum)
	// d.加入到原文中
	contentByte := append([]byte(contentStr), padByte...)
	// 4.加密
	cryptByte := make([]byte, len(contentByte))
	blockMode.CryptBlocks(cryptByte, contentByte)
	// 5.base编码返回
	return base64.StdEncoding.EncodeToString(cryptByte), nil
}

// AesDecrypt AES解密
func AesDecrypt(cryptStr string, key string, iv string) (string, error) {
	// 1.创建AES加密对象
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", nil
	}
	// 2.CBC加密模式
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))

	// 3.Base64解码
	cryptByte, err := base64.StdEncoding.DecodeString(cryptStr)
	if err != nil {
		return "", err
	}
	// 4.解码
	contentByte := make([]byte, len(cryptByte))
	blockMode.CryptBlocks(contentByte, cryptByte)
	// 5.去除填充字节
	length := len(contentByte)
	unPadding := int(contentByte[length-1])
	return string(contentByte[:length-unPadding]), nil
}
