package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

type Aes struct {
	Key string
	Iv  string
}

//Encode 开始加密
func (a *Aes) Encode(data string) (string, error) {
	_data := []byte(data)
	_key := []byte(a.Key)
	_iv := []byte(a.Iv)

	_data = a.PKCS7Padding(_data)
	block, err := aes.NewCipher(_key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, _iv)
	mode.CryptBlocks(_data, _data)
	return base64.StdEncoding.EncodeToString(_data), nil
}

//Decode 开始解密
func (a *Aes) Decode(data string) (str string, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = errors.New(fmt.Sprintf("%v", e))
		}
	}()
	_data, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	_key := []byte(a.Key)
	_iv := []byte(a.Iv)

	block, err := aes.NewCipher(_key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, _iv)
	mode.CryptBlocks(_data, _data)
	_data = a.PKCS7UnPadding(_data)

	return string(_data), nil
}
func (a *Aes) PKCS7Padding(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}
func (a *Aes) PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
