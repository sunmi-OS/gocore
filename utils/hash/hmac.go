package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

func HmacSha1(value, secret string) (string, error) {
	return Hmac(sha1.New, value, secret)
}

func HmacSha224(value, secret string) (string, error) {
	return Hmac(sha256.New224, value, secret)
}

func HmacSha256(value, secret string) (string, error) {
	return Hmac(sha256.New, value, secret)
}

func HmacSha384(value, secret string) (string, error) {
	return Hmac(sha512.New384, value, secret)
}

func HmacSha512(value, secret string) (string, error) {
	return Hmac(sha512.New, value, secret)
}

func HmacMD5(value, secret string) (string, error) {
	return Hmac(md5.New, value, secret)
}

func Hmac(h func() hash.Hash, value, key string) (string, error) {
	mac := hmac.New(h, []byte(key))
	_, err := mac.Write([]byte(value))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}
