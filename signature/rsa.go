package signature

import (
    "encoding/pem"
    "crypto/x509"
    "crypto/rsa"
    "crypto"
    "encoding/base64"
    "crypto/rand"
    "crypto/sha1"
    "crypto/sha256"
)

func FormatPrivateKey(privateKey string) string {
    tempKey := []byte(privateKey)
    length := len(tempKey)
    formatPriKey := "-----BEGIN RSA PRIVATE KEY-----\n"
    tail := make([]byte, length)
    for i := 0; i < length; i ++ {
        if (i + 1) % 64 == 0 {
            head := tempKey[i-63:i+1]
            tail = tempKey[i+1:]
            formatPriKey += string(head) + "\n"
        }
    }
    formatPriKey += string(tail) + "\n"
    formatPriKey += "-----END RSA PRIVATE KEY-----\n"
    return formatPriKey
}

func FastFormatPrivateKey(privateKey string) []byte {
    tempKey := []byte(privateKey)
    length := len(tempKey)
    page := length / 64
    formatPriKey := "-----BEGIN RSA PRIVATE KEY-----\n"
    for i := 0; i < page; i ++ {
        formatPriKey += string(tempKey[i*64:(i+1)*64]) + "\n"
    }
    formatPriKey += string(tempKey[page*64:]) + "\n"
    formatPriKey += "-----END RSA PRIVATE KEY-----\n"
    return []byte(formatPriKey)
}

func SignSha1WithRsaPKCS8(data string, privateKey []byte) (string, error) {
    block, _ := pem.Decode(privateKey)
    if block == nil {
        panic("RsaSign PrivateKey Error")
    }
    p, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        panic(err.Error())
    }
    pri := p.(*rsa.PrivateKey)

    sha1Hash := sha1.New()
    s_data := []byte(data)
    sha1Hash.Write(s_data)
    hashed := sha1Hash.Sum(nil)

    signByte, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA1, hashed)
    sign := base64.StdEncoding.EncodeToString(signByte)
    return string(sign), err
}

func VerifySignSha1WithRsa(data string, signData string, publicKey string) error {
    sign, err := base64.StdEncoding.DecodeString(signData)
    if err != nil {
        return err
    }
    public, _ := base64.StdEncoding.DecodeString(publicKey)
    pub, err := x509.ParsePKIXPublicKey(public)
    if err != nil {
        return err
    }
    hash := sha1.New()
    hash.Write([]byte(data))
    return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

func VerifySignSha256WithRsa(data string, signData string, publicKey string) error {
    sign, err := base64.StdEncoding.DecodeString(signData)
    if err != nil {
        return err
    }
    public, _ := base64.StdEncoding.DecodeString(publicKey)
    pub, err := x509.ParsePKIXPublicKey(public)
    if err != nil {
        return err
    }
    hash := sha256.New()
    hash.Write([]byte(data))

    return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hash.Sum(nil), sign)
}

func SignSha256WithRsaPKCS8(data string, privateKey []byte) (string, error) {
    block, _ := pem.Decode(privateKey)
    if block == nil {
        panic("RsaSign PrivateKey Error")
    }
    p, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        panic(err.Error())
    }
    pri := p.(*rsa.PrivateKey)

    sha256Hash := sha256.New()
    s_data := []byte(data)
    sha256Hash.Write(s_data)
    hashed := sha256Hash.Sum(nil)

    signByte, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, hashed)
    sign := base64.StdEncoding.EncodeToString(signByte)
    return string(sign), err
}

func SignSha1WithRsaPKCS1(data string, privateKey []byte) (string, error) {
    block, _ := pem.Decode(privateKey)
    if block == nil {
        panic("RsaSign PrivateKey Error")
    }
    pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        panic(err.Error())
    }
    sha1Hash := sha1.New()
    s_data := []byte(data)
    sha1Hash.Write(s_data)
    hashed := sha1Hash.Sum(nil)

    signByte, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA1, hashed)
    sign := base64.StdEncoding.EncodeToString(signByte)
    return string(sign), err
}

func SignSha256WithRsaPKCS1(data string, privateKey []byte) (string, error) {
    block, _ := pem.Decode(privateKey)
    if block == nil {
        panic("RsaSign PrivateKey Error")
    }
    pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        panic(err.Error())
    }
    sha256Hash := sha256.New()
    s_data := []byte(data)
    sha256Hash.Write(s_data)
    hashed := sha256Hash.Sum(nil)

    signByte, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, hashed)
    sign := base64.StdEncoding.EncodeToString(signByte)
    return string(sign), err
}
