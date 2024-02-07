package rsacrypto

import (
	"crypto/rsa"
	"encoding/hex"
	"testing"
)

var privkeyPem = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAx5Hg78/Gwr4dtQ3ZehrQVet+RO5/k7Xm1RmNly0sCdPhcsaw
a8B0o7L7L11k2NP9ypVn2CDRHx8znhrtzSfBzh/HLCA9XLhG6jpABDsgam15hSvw
bVgZNvTRXfDZvZXpK3fP2tH/z9ncxeQxrVJrOg/zf8GskMscTOQ1VdxiKQ2uZPuo
BxHYe+UDkhWpYfZAcDSwUvO0JWUAdvteJT3o+kBcA3QywVAohfnzepej+QAZXFNf
3Hs+PfoxOLejQpsdZ8KKonmtbrwUIuyDLVdWYtsEDA8g64rmIwKxPlNpFTyRXgMA
OEkrEpiCRb7HnnOE6XAxJlou2z6HfB41QoO9+wIDAQABAoIBAHcLQKcsRL7j0yqu
ET0yA3ZNHCwYLEe7KO+S54/3JR7TodbqSFBuI+WGHSmax05D3k7aonAc20F6RjsY
iyNmhMfk0tUyggft8HdFuewMLQDvPp6+oBNJivjqPn2P7wKVCtqgBH/d5n9g0L3G
qg9ea5Hd8/0QVVSlo8MGGf6WkIM1lOrsPi6iBicweFEIeBRy9RBUv1KKJVThpZhZ
SmCM7Vy9nLVxegv+5tbbYACnQC4CozVDPKmz0UFe+P7TQ07K0g/SNAF15Hl+lnMJ
HdrecqZac7qGEJnrCI477o8gitlBo579OZeZkJd0JM2GrUjzcCOJWAHTKhDbs6Di
PlUzU2ECgYEA7zR74nrWA5YrMklUOk9ssjuiv1gA9tN+F54JZN6AKp3nMguHlO2t
ViM0yArQBzTmresJ6hMwbemEU8LLFfaW5aODVZ8oxgm6gfbhxNnAYMwI+k/Itx0a
ooj5OK8y6rpcgHJzlAwOxPqn9VJ3Jze6ll1KPpeA/A9hhgFelAWu/OsCgYEA1ZT8
DgHyrgzjZagR+sI7pRBxfUKUO4xi/LcwQhi3Z71RhoSSS+vzbaA2lyOw8vtPj6gE
K8nuTAruuLfpVEdmIbu4EtoI+Bw9K0w7Cg2mXVGHvrJNlUTRY1wwnNMAVzWeIwDC
OnvnDGRodv0lLRxs0XPjtLuNx7wnrnk8l9t2vzECgYEA7T3JjNM1hXMfvo0Z24dA
j/kzrcDzm9ogmf3k5UUEKsBXN7xVqTCdlOvwAmMu9abTDzUorR6BDtHmq0hsMYlT
Gci1jmr/foLRluqr+pfZBGf4k4Ij2PElpIRjYYPp5QIWklJxLSlUUKslf9tdT+km
xtEZvMB4bgY3PDgJfJeyeScCgYEAyk5Vtfr4YQ7KMldRuIFUt9Rse2aePA2NEa1/
Y4w/5V65Iz7dyFZV/Sf9rYncKTwMr5lJYiTiuFq+pm9l7zO2NQu3nvux9TniYunR
HoOxasE4YFRKErLd10zSqyleMD0UbjlgwL7uKpnNLbA5D5LWLEumi2IAOQorWCN0
Vq9FunECgYAaZAV4POfwGjJsJsBlAyNeddSytANVwFoioOQ9Wo6/E66MioAbZU3/
fUEzmjo90ainHi8BLPiL9oEM9chn7Rq2kfUK4XDlo3/3k3O8rsQ/3DDEfD4UagAZ
FqwC148L2uFYVuCJ3huG+Ab0q0uQ1xlrY4rV1aRl5GLJwVDQFq3eEg==
-----END RSA PRIVATE KEY-----`

var pubkeyPem = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAx5Hg78/Gwr4dtQ3ZehrQ
Vet+RO5/k7Xm1RmNly0sCdPhcsawa8B0o7L7L11k2NP9ypVn2CDRHx8znhrtzSfB
zh/HLCA9XLhG6jpABDsgam15hSvwbVgZNvTRXfDZvZXpK3fP2tH/z9ncxeQxrVJr
Og/zf8GskMscTOQ1VdxiKQ2uZPuoBxHYe+UDkhWpYfZAcDSwUvO0JWUAdvteJT3o
+kBcA3QywVAohfnzepej+QAZXFNf3Hs+PfoxOLejQpsdZ8KKonmtbrwUIuyDLVdW
YtsEDA8g64rmIwKxPlNpFTyRXgMAOEkrEpiCRb7HnnOE6XAxJlou2z6HfB41QoO9
+wIDAQAB
-----END RSA PUBLIC KEY-----`

var modulusHex = "c791e0efcfc6c2be1db50dd97a1ad055eb7e44ee7f93b5e6d5198d972d2c09d3e172c6b06bc074a3b2fb2f5d64d8d3fdca9567d820d11f1f339e1aedcd27c1ce1fc72c203d5cb846ea3a40043b206a6d79852bf06d581936f4d15df0d9bd95e92b77cfdad1ffcfd9dcc5e431ad526b3a0ff37fc1ac90cb1c4ce43555dc62290dae64fba80711d87be5039215a961f6407034b052f3b425650076fb5e253de8fa405c037432c1502885f9f37a97a3f900195c535fdc7b3e3dfa3138b7a3429b1d67c28aa279ad6ebc1422ec832d575662db040c0f20eb8ae62302b13e5369153c915e030038492b12988245bec79e7384e97031265a2edb3e877c1e354283bdfb"
var pubExponent = "10001"
var privExponent = "770b40a72c44bee3d32aae113d3203764d1c2c182c47bb28ef92e78ff7251ed3a1d6ea48506e23e5861d299ac74e43de4edaa2701cdb417a463b188b236684c7e4d2d5328207edf07745b9ec0c2d00ef3e9ebea013498af8ea3e7d8fef02950adaa0047fdde67f60d0bdc6aa0f5e6b91ddf3fd105554a5a3c30619fe9690833594eaec3e2ea2062730785108781472f51054bf528a2554e1a598594a608ced5cbd9cb5717a0bfee6d6db6000a7402e02a335433ca9b3d1415ef8fed3434ecad20fd2340175e4797e9673091ddade72a65a73ba861099eb088e3bee8f208ad941a39efd39979990977424cd86ad48f37023895801d32a10dbb3a0e23e55335361"

func TestDecodePubKeyPem(t *testing.T) {
	pub, err := DecodePubKeyPem(pubkeyPem)
	if err != nil {
		t.Error(err)
	}
	plaintext := "hello world"
	ciphertext := "02ceb87bc708d90f7a5cab39c0a8f55945a3e10a6257568760f429afccba83f791271a6a9a9465980b130fe27b033306506b95c9e792a8c2e4c69f18265646ccdd456f298bd89f808001230a6883b2a3f13f143f8f8f3c1ecc51b9e7d0e38daef5d383cd8248e80ea49732f1aa91840e5a07f1c57ebdf0d96611050fa841cdfa02c7275dad403b026221b8b0007100e0d928d2b05d42bb6d7862be341c4ada2c086676ca2954dbff8123f3627f681cbae13abf43ceeebd16fc4a155e8c3b04223bcc094b40985abdaec3bc56fe193a1926fd0e9f35116ca31a3bd6437f56141023045a30c66bbb13c1cd2e0e9e2957a3778884bcded5b0bd415fbb605fbacb4c"
	b, _ := hex.DecodeString(ciphertext)
	result, err := DecryptPKCS1WithPubkey(pub, b)
	if err != nil {
		t.Error(err)
	}
	if string(result) != plaintext {
		t.Error("result is not equal to plaintext")
	}
}

func TestDecodePrivKeyPem(t *testing.T) {
	priv, err := DecodePrivKeyPem(privkeyPem)
	if err != nil {
		t.Error(err)
	}
	plaintext := "hello world"
	b, err := EncryptPKCS1WithPrivkey(priv, []byte(plaintext))
	if err != nil {
		t.Error(err)
	}
	result, err := DecryptPKCS1WithPubkey(&priv.PublicKey, b)
	if err != nil {
		t.Error(err)
	}
	if string(result) != plaintext {
		t.Error("result is not equal to plaintext")
	}
}

func TestEncryptPKCS1WithPrivkey(t *testing.T) {
	priv := decodePrivKey(privkeyPem)
	plainText := "hello world"
	b, err := EncryptPKCS1WithPrivkey(priv, []byte(plainText))
	if err != nil {
		t.Error(err)
	}
	result, err := DecryptPKCS1WithPubkey(&priv.PublicKey, b)
	if err != nil {
		t.Error(err)
	}
	if string(result) != plainText {
		t.Error("result is not equal to plainText")
	}
}

func TestDecryptPKCS1WithPubkey(t *testing.T) {
	priv := decodePrivKey(privkeyPem)
	plaintext := "hello world"
	b, err := EncryptPKCS1WithPrivkey(priv, []byte(plaintext))
	if err != nil {
		t.Error(err)
	}
	result, err := DecryptPKCS1WithPubkey(&priv.PublicKey, b)
	if err != nil {
		t.Error(err)
	}
	if string(result) != plaintext {
		t.Error("result is not equal to plaintext")
	}
}

func TestDecodePubKeyModulus(t *testing.T) {
	pub, err := DecodePubKeyModulus(modulusHex, pubExponent)
	if err != nil {
		t.Error(err)
	}
	plaintext := "hello world"
	ciphertext := "02ceb87bc708d90f7a5cab39c0a8f55945a3e10a6257568760f429afccba83f791271a6a9a9465980b130fe27b033306506b95c9e792a8c2e4c69f18265646ccdd456f298bd89f808001230a6883b2a3f13f143f8f8f3c1ecc51b9e7d0e38daef5d383cd8248e80ea49732f1aa91840e5a07f1c57ebdf0d96611050fa841cdfa02c7275dad403b026221b8b0007100e0d928d2b05d42bb6d7862be341c4ada2c086676ca2954dbff8123f3627f681cbae13abf43ceeebd16fc4a155e8c3b04223bcc094b40985abdaec3bc56fe193a1926fd0e9f35116ca31a3bd6437f56141023045a30c66bbb13c1cd2e0e9e2957a3778884bcded5b0bd415fbb605fbacb4c"
	b, _ := hex.DecodeString(ciphertext)
	result, err := DecryptPKCS1WithPubkey(pub, b)
	if err != nil {
		t.Error(err)
	}
	if string(result) != plaintext {
		t.Error("result is not equal to plaintext")
	}
}

func TestDecodePrivKeyModulus(t *testing.T) {
	priv, err := DecodePrivKeyModulus(modulusHex, pubExponent, privExponent)
	if err != nil {
		t.Error(err)
	}
	plaintext := "hello world"
	b, err := EncryptPKCS1WithPrivkey(priv, []byte(plaintext))
	if err != nil {
		t.Error(err)
	}
	result, err := DecryptPKCS1WithPubkey(&priv.PublicKey, b)
	if err != nil {
		t.Error(err)
	}
	if string(result) != plaintext {
		t.Error("result is not equal to plaintext")
	}
}

func decodePrivKey(s string) *rsa.PrivateKey {
	priv, err := DecodePrivKeyPem(s)
	if err != nil {
		panic(err)
	}
	return priv
}
