package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"unsafe"
)

// EncryptCBC DES/CBC/PKCS5Padding   加密
func EncryptCBC(msg, key, iv string) (string, error) {
	origData := []byte(msg)
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypted := make([]byte, len(origData))
	// pre-verification to prevent panic
	err = checkBlock(block, crypted, origData)
	if err != nil {
		return "", err
	}
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return string(crypted), nil
}

// DecryptCBC DES/CBC/PKCS5Padding  解密
func DecryptCBC(msg, key, iv string) (string, error) {
	crypted := []byte(msg)
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	origData := make([]byte, len(crypted))
	// pre-verification to prevent panic
	err = checkBlock(block, origData, crypted)
	if err != nil {
		return "", err
	}
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return "", err
	}
	// origData = ZeroUnPadding(origData)
	return string(origData), nil
}

// EncryptECB DES/ECB/PKCS5Padding
func EncryptECB(msg, key string) (string, error) {
	origData := []byte(msg)
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	origData = PKCS5Padding(origData, bs)
	if len(origData)%bs != 0 {
		return "", errors.New("invalid key size")
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

// DecryptECB DES/ECB/PKCS5Padding
func DecryptECB(msg, key string) (string, error) {
	crypted := []byte(msg)
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return "", errors.New("invalid key size or data size")
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
	err = checkBlock(block, crypted, origData)
	if err != nil {
		return "", err
	}
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
	err = checkBlock(block, origData, crypted)
	if err != nil {
		return "", err
	}
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

func checkBlock(b cipher.Block, dst, src []byte) error {
	if len(src)%b.BlockSize() != 0 {
		return errors.New("input not full blocks")
	}
	if len(dst) < len(src) {
		return errors.New("output smaller than input")
	}
	if inexactOverlap(dst[:len(src)], src) {
		return errors.New("invalid buffer overlap")
	}
	return nil
}

// InexactOverlap reports whether x and y share memory at any non-corresponding
// index. The memory beyond the slice length is ignored. Note that x and y can
// have different lengths and still not have any inexact overlap.
//
// InexactOverlap can be used to implement the requirements of the crypto/cipher
// AEAD, Block, BlockMode and Stream interfaces.
func inexactOverlap(x, y []byte) bool {
	if len(x) == 0 || len(y) == 0 || &x[0] == &y[0] {
		return false
	}
	return anyOverlap(x, y)
}

// corresponding) index. The memory beyond the slice length is ignored.
func anyOverlap(x, y []byte) bool {
	return len(x) > 0 && len(y) > 0 &&
		uintptr(unsafe.Pointer(&x[0])) <= uintptr(unsafe.Pointer(&y[len(y)-1])) &&
		uintptr(unsafe.Pointer(&y[0])) <= uintptr(unsafe.Pointer(&x[len(x)-1]))
}
