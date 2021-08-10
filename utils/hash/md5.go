package hash

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// MD5 加密字符串
func MD5(plainText string) string {
	h := md5.New()
	_, err := h.Write([]byte(plainText))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

// MD5File 计算文件的md5，适用于本地文件计算
func MD5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(md5hash.Sum(nil)), nil
}
