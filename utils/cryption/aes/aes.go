package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

// EncryptUseCBCWithDefaultProtocol Encrypt using random iv parameter and cbc mode
// Never panic, only possible to return an error
func EncryptUseCBCWithDefaultProtocol(plainText, key []byte) ([]byte, error) {
	iv := make([]byte, 16)
	// random iv param
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}
	cipherText, err := EncryptUseCBC(plainText, key, iv)
	if err != nil {
		return nil, err
	}
	// Put the iv parameter in the head of the cipher text
	result := append(iv, cipherText...)
	return result, err
}

// EncryptUseCBC Encrypt using cbc mode
// When iv length does not equal block size, it will panic
func EncryptUseCBC(plainText, key, iv []byte) ([]byte, error) {
	blockKey, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := blockKey.BlockSize()
	// do padding
	fixedPlainText := PKCS5Padding(plainText, blockSize)
	encryptTool := cipher.NewCBCEncrypter(blockKey, iv)
	cipherText := make([]byte, len(fixedPlainText))
	// do final
	encryptTool.CryptBlocks(cipherText, fixedPlainText)
	return cipherText, nil
}

// DecryptUseCBC Decrypt using cbc mode
// There are two kinds of panic that may occur:
// 1. When iv length do not equal block size
// 2. When key does not match the cipher text, and it always happens when do unpadding
func DecryptUseCBC(cipherText, key []byte, iv []byte) ([]byte, error) {
	blockKey, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := blockKey.BlockSize()
	if len(cipherText)%blockSize != 0 {
		return nil, errors.New("cipher text is not an integral multiple of the block size")
	}
	decryptTool := cipher.NewCBCDecrypter(blockKey, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	decryptTool.CryptBlocks(cipherText, cipherText)
	origData, err := PKCS5UnPadding(cipherText)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

// DecryptUseCBCWithDefaultProtocol Decrypt using given iv parameter and cbc mode
// When key does not match the cipher text, it will panic
func DecryptUseCBCWithDefaultProtocol(cipherText, key []byte) ([]byte, error) {
	if len(cipherText) < 16 {
		return nil, errors.New("decrypt excepted iv parameter")
	}
	plainText, err := DecryptUseCBC(cipherText[16:], key, cipherText[:16])
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func getKey(key string) []byte {
	keyLen := len(key)
	if keyLen < 16 {
		panic("res key 长度不能小于16")
	}
	arrKey := []byte(key)
	if keyLen >= 32 {
		// 取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		// 取前24个字节
		return arrKey[:24]
	}
	// 取前16个字节
	return arrKey[:16]
}

// Base64UrlSafeEncode Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
func Base64UrlSafeEncode(source []byte) string {
	byteArr := base64.StdEncoding.EncodeToString(source)
	safeUrl := strings.Replace(byteArr, "/", "_", -1)
	safeUrl = strings.Replace(safeUrl, "+", "-", -1)
	safeUrl = strings.Replace(safeUrl, "=", "", -1)
	return safeUrl
}

func AesDecrypt(msg, k string) (string, error) {
	key := getKey(k)
	crypted, _ := base64.StdEncoding.DecodeString(msg)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	err = checkBlock(block, origData, crypted)
	if err != nil {
		return "", err
	}
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}

func AesEncrypt(src, k string) (string, error) {
	key := getKey(k)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	if src == "" {
		return "", errors.New("encrypt data is empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	err = checkBlock(block, crypted, content)
	if err != nil {
		return "", err
	}
	ecb.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func checkBlock(b cipher.Block, dst, src []byte) error {
	if len(src)%b.BlockSize() != 0 {
		return errors.New("input data is incomplete")
	}
	if len(dst) < len(src) {
		return errors.New("crypto/cipher: output smaller than input")
	}
	return nil
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

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func EncryptUseCTRNoPadding(plainText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// padding
	fixedPlainText := NoPadding(plainText, block.BlockSize())
	// mode
	blockMode := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(fixedPlainText))
	// do final
	blockMode.XORKeyStream(cipherText, fixedPlainText)
	return cipherText, nil
}

func DecryptUseCTRNoPadding(cipherText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// mode
	blockMode := cipher.NewCTR(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.XORKeyStream(origData, cipherText)
	return NoUnPadding(origData), nil
}

func NoPadding(cipherText []byte, blockSize int) []byte {
	return cipherText
}

func NoUnPadding(origData []byte) []byte {
	return origData
}
