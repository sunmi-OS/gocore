package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

func DesEncrypt(msg string, key string, iv string) (string, error) {

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

func DesDecrypt(msg string, key string, iv string) (string, error) {

	crypted := []byte(msg)
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return string(origData), nil
}

/*
DES/ECB/PKCS5Padding   加密
*/
func DesEncryptECB(msg string, key string) (string, error) {

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

/*
DES/ECB/PKCS5Padding   解密
*/
func DesDecryptECB(msg string, key string) (string, error) {

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
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
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

// 3DES加密
func TripleDesEncrypt(msg string, key string, iv string) (string, error) {

	origData := []byte(msg)

	block, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	//blockMode := cipher.NewCBCEncrypter(block, key[:8])
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return string(crypted), nil
}

// 3DES解密
func TripleDesDecrypt(msg string, key string, iv string) (string, error) {

	crypted := []byte(msg)
	block, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		return "", err
	}
	//blockMode := cipher.NewCBCDecrypter(block, key[:8])
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return string(origData), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
