package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
	"password-self-service/pkg/config"
)

/*
生成1024位的RSA私钥：
openssl genrsa -out private_key.pem 1024
根据私钥生成公钥：
openssl rsa -in private_key.pem -pubout -out public_key.pem
*/

var (
	publicKey  []byte
	privateKey []byte
)

// ReadFile 读取文件
func ReadFile(keyFile string) ([]byte, error) {
	if f, err := os.Open(keyFile); err != nil {
		return nil, err
	} else {
		content := make([]byte, 4096)
		if n, err := f.Read(content); err != nil {
			return nil, err
		} else {
			return content[:n], err
		}
	}
}

// ReadRSAKey 读取公私钥文件
func ReadRSAKey() {
	publicKey, _ = ReadFile(config.Setting.Server.RSAPublicKey)
	privateKey, _ = ReadFile(config.Setting.Server.RSAPrivateKey)
}

// Encrypt Rsa加密
func Encrypt(origData []byte) (string, error) {
	// 解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}

	// 解析公钥，目前数字证书一般都是基于ITU（国际电信联盟）指定的x.509标准
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)

	// 加密明文
	encryptText, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origData)

	// 转为base64
	return base64.StdEncoding.EncodeToString(encryptText), err
}

// Decrypt Rsa解密
func Decrypt(cipherText string) (string, error) {
	// base64解码
	encryptText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	// 解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error")
	}

	// 解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 解密密文
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, priv, encryptText)

	return string(plainText), err
}
