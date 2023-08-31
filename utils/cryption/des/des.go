package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

// EncryptCBC DES/CBC/PKCS5Padding   加密
func EncryptCBC(msg string, key string, iv string) (string, error) {

	origData := []byte(msg)
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return string(crypted), nil
}

// DecryptCBC DES/CBC/PKCS5Padding  解密
func DecryptCBC(msg string, key string, iv string) (string, error) {
	crypted := []byte(msg)
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return "", err
	}
	// origData = ZeroUnPadding(origData)
	return string(origData), nil
}

// EncryptECB DES/ECB/PKCS5Padding   加密
func EncryptECB(msg string, key string) (string, error) {

	origData := []byte(msg)

	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	origData = PKCS5Padding(origData, bs)
	if len(origData)%bs != 0 {
		return "", err
	}
	crypted := make([]byte, len(origData))
	dst := crypted
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}

	return string(crypted), nil
}

// DecryptECB DES/ECB/PKCS5Padding   解密
func DecryptECB(msg string, key string) (string, error) {

	crypted := []byte(msg)

	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return "", err
	}
	origData := make([]byte, len(crypted))
	dst := origData
	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}
	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}

// TripleEncrypt 3DES加密
func TripleEncrypt(msg string, key string, iv string) (string, error) {
	origData := []byte(msg)
	block, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return string(crypted), nil
}

// TripleDecrypt 3DES解密
func TripleDecrypt(msg string, key string, iv string) (string, error) {
	crypted := []byte(msg)
	block, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("decryption failure")
	}
	unpadding := int(origData[length-1])
	if unpadding > length {
		return nil, errors.New("decryption failure")
	}
	return origData[:(length - unpadding)], nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}
