package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func Sha1(text string) (string, error) {
	sha := sha1.New()
	_, err := sha.Write([]byte(text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}

func Sha224(text string) (string, error) {
	sha := sha256.New224()
	_, err := sha.Write([]byte(text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}

func Sha256(text string) (string, error) {
	sha := sha256.New()
	_, err := sha.Write([]byte(text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}

func Sha384(text string) (string, error) {
	sha := sha512.New384()
	_, err := sha.Write([]byte(text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}

func Sha512(text string) (string, error) {
	sha := sha512.New()
	_, err := sha.Write([]byte(text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}
