package enc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"log"
)

func AESEncrypt(orign, key, iv []byte) ([]byte, error) {
	sKey := make([]byte, 16)
	a := md5.Sum(key)
	copy(sKey, a[:16])
	b, err := aes.NewCipher(sKey)
	if err != nil {
		return nil, err
	}
	bSize := b.BlockSize()
	bMode := cipher.NewCBCEncrypter(b, iv[:bSize])
	orign = PKCS7Padding(orign, bSize)
	crypt := make([]byte, len(orign))
	bMode.CryptBlocks(crypt, orign)
	return crypt, nil
}

func AESDecrypt(crypt, key, iv []byte) ([]byte, error) {
	sKey := make([]byte, 16)
	a := md5.Sum(key)
	copy(sKey, a[:16])
	b, err := aes.NewCipher(sKey)
	if err != nil {
		return nil, err
	}
	bSize := b.BlockSize()
	bMode := cipher.NewCBCDecrypter(b, iv[:bSize])
	orign := make([]byte, len(crypt))
	log.Println(len(crypt))
	bMode.CryptBlocks(orign, crypt)
	orign = PKCS7UnPadding(orign)
	return orign, nil
}

func PKCS7Padding(orign []byte, blockSize int) []byte {
	padding := blockSize - len(orign)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(orign, padtext...)
}

func PKCS7UnPadding(orign []byte) []byte {
	length := len(orign)
	unpadding := int(orign[length-1])
	return orign[:(length - unpadding)]
}
