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

func HmacSHA1(secret, params string) (string, error) {
	return Hmac(sha1.New, []byte(secret), []byte(params))
}

func HmacSHA224(secret, params string) (string, error) {
	return Hmac(sha256.New224, []byte(secret), []byte(params))
}

func HmacSHA256(secret, params string) (string, error) {
	return Hmac(sha256.New, []byte(secret), []byte(params))
}

func HmacSha384(secret, params string) (string, error) {
	return Hmac(sha512.New384, []byte(secret), []byte(params))
}

func HmacSHA512Sign(secret, params string) (string, error) {
	return Hmac(sha512.New, []byte(secret), []byte(params))
}

func HmacMD5(secret, params string) (string, error) {
	return Hmac(md5.New, []byte(secret), []byte(params))
}

func Hmac(h func() hash.Hash, key, params []byte) (string, error) {
	mac := hmac.New(h, key)
	_, err := mac.Write(params)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}
