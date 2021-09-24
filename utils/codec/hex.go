package codec

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

// Base64Encode 字符串转64进制
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode 64进制转字符串
func Base64Decode(str string) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(sDec), nil
}

// Base32Encode 字符串转32进制
func Base32Encode(str string) string {
	return base32.StdEncoding.EncodeToString([]byte(str))
}

// Base32Decode 32进制转字符串
func Base32Decode(str string) (string, error) {
	sDec, err := base32.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(sDec), nil
}

// HexEncode 字符串转16进制
func HexEncode(str string) string {
	return hex.EncodeToString([]byte(str))
}

// HexDecode 16进制转字符串
func HexDecode(str string) (string, error) {
	sDec, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(sDec), nil
}
