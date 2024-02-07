package rsacrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"
	"io"
	"math/big"
	"regexp"
)

// EncryptOAEP encrypts the given message with RSA-OAEP.
//
// OAEP is parameterised by a hash function that is used as a random oracle.
// Encryption and decryption of a given message must use the same hash function
// and sha256.New() is a reasonable choice.
//
// The random parameter is used as a source of entropy to ensure that
// encrypting the same message twice doesn't result in the same ciphertext.
// Most applications should use [crypto/rand.Reader] as random.
//
// The label parameter may contain arbitrary data that will not be encrypted,
// but which gives important context to the message. For example, if a given
// public key is used to encrypt two types of messages then distinct label
// values could be used to ensure that a ciphertext for one purpose cannot be
// used for another by an attacker. If not required it can be empty.
//
// The message must be no longer than the length of the public modulus minus
// twice the hash length, minus a further 2.
func EncryptOAEP(hash hash.Hash, random io.Reader, pub *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	return rsa.EncryptOAEP(hash, random, pub, msg, label)
}

// EncryptPKCS1v15 encrypts the given message with RSA and the padding
// scheme from PKCS #1 v1.5.  The message must be no longer than the
// length of the public modulus minus 11 bytes.
//
// The random parameter is used as a source of entropy to ensure that
// encrypting the same message twice doesn't result in the same
// ciphertext. Most applications should use [crypto/rand.Reader]
// as random. Note that the returned ciphertext does not depend
// deterministically on the bytes read from random, and may change
// between calls and/or between versions.
//
// WARNING: use of this function to encrypt plaintexts other than
// session keys is dangerous. Use RSA OAEP in new protocols.
func EncryptPKCS1v15(random io.Reader, pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(random, pub, msg)
}

// DecryptOAEP decrypts ciphertext using RSA-OAEP.
//
// OAEP is parameterised by a hash function that is used as a random oracle.
// Encryption and decryption of a given message must use the same hash function
// and sha256.New() is a reasonable choice.
//
// The random parameter is legacy and ignored, and it can be nil.
//
// The label parameter must match the value given when encrypting. See
// EncryptOAEP for details.
func DecryptOAEP(hash hash.Hash, random io.Reader, priv *rsa.PrivateKey, ciphertext []byte, label []byte) ([]byte, error) {
	return rsa.DecryptOAEP(hash, random, priv, ciphertext, label)
}

// DecryptPKCS1v15 decrypts a plaintext using RSA and the padding scheme from PKCS #1 v1.5.
// The random parameter is legacy and ignored, and it can be nil.
//
// Note that whether this function returns an error or not discloses secret
// information. If an attacker can cause this function to run repeatedly and
// learn whether each instance returned an error then they can decrypt and
// forge signatures as if they had the private key. See
// DecryptPKCS1v15SessionKey for a way of solving this problem.
func DecryptPKCS1v15(random io.Reader, priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(random, priv, ciphertext)
}

// DecryptPKCS1v15SessionKey decrypts a session key using RSA and the padding
// scheme from PKCS #1 v1.5. The random parameter is legacy and ignored, and it
// can be nil.
func DecryptPKCS1v15SessionKey(random io.Reader, priv *rsa.PrivateKey, ciphertext []byte, key []byte) error {
	return rsa.DecryptPKCS1v15SessionKey(random, priv, ciphertext, key)
}

// SignPKCS1v15 calculates the signature of hashed using
// RSASSA-PKCS1-V1_5-SIGN from RSA PKCS #1 v1.5.  Note that hashed must
// be the result of hashing the input message using the given hash
// function. If hash is zero, hashed is signed directly. This isn't
// advisable except for interoperability.
//
// The random parameter is legacy and ignored, and it can be nil.
//
// This function is deterministic. Thus, if the set of possible
// messages is small, an attacker may be able to build a map from
// messages to signatures and identify the signed messages. As ever,
// signatures provide authenticity, not confidentiality.
func SignPKCS1v15(random io.Reader, priv *rsa.PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(random, priv, hash, hashed)
}

// SignPSS calculates the signature of digest using PSS.
//
// digest must be the result of hashing the input message using the given hash
// function. The opts argument may be nil, in which case sensible defaults are
// used. If opts.Hash is set, it overrides hash.
//
// The signature is randomized depending on the message, key, and salt size,
// using bytes from rand. Most applications should use [crypto/rand.Reader] as
// rand.
func SignPSS(rand io.Reader, priv *rsa.PrivateKey, hash crypto.Hash, digest []byte, opts *rsa.PSSOptions) ([]byte, error) {
	return rsa.SignPSS(rand, priv, hash, digest, opts)
}

// VerifyPKCS1v15 verifies an RSA PKCS #1 v1.5 signature.
// hashed is the result of hashing the input message using the given hash
// function and sig is the signature. A valid signature is indicated by
// returning a nil error. If hash is zero then hashed is used directly. This
// isn't advisable except for interoperability.
func VerifyPKCS1v15(pub *rsa.PublicKey, hash crypto.Hash, hashed []byte, sig []byte) error {
	return rsa.VerifyPKCS1v15(pub, hash, hashed, sig)
}

// VerifyPSS verifies a PSS signature.
//
// A valid signature is indicated by returning a nil error. digest must be the
// result of hashing the input message using the given hash function. The opts
// argument may be nil, in which case sensible defaults are used. opts.Hash is
// ignored.
func VerifyPSS(pub *rsa.PublicKey, hash crypto.Hash, digest []byte, sig []byte, opts *rsa.PSSOptions) error {
	return rsa.VerifyPSS(pub, hash, digest, sig, opts)
}

// EncryptPKCS1WithPrivkey encryption with private key
//
// In normal security practice, we do not use private keys to encrypt data, but public keys,
// This function is provided because such irregularities do exist in practice.
func EncryptPKCS1WithPrivkey(priv *rsa.PrivateKey, msg []byte) ([]byte, error) {
	tLen := len(msg)
	k := (priv.N.BitLen() + 7) / 8
	if k < tLen+11 {
		return nil, errors.New("rsacrypto: message too long for rsa key size")
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], msg)
	m := new(big.Int).SetBytes(em)
	c, err := decrypt(rand.Reader, priv, m)
	if err != nil {
		return nil, err
	}
	copyWithLeftPad(em, c.Bytes())
	return em, nil
}

// DecryptPKCS1WithPubkey Decryption with public key
//
// In normal security practice, we do not use public keys to decrypt data, but private keys,
// This function is provided because such irregularities do exist in practice.
func DecryptPKCS1WithPubkey(pub *rsa.PublicKey, ciphertext []byte) ([]byte, error) {
	if err := checkPub(pub); err != nil {
		return nil, err
	}
	k := (pub.N.BitLen() + 7) / 8
	if k != len(ciphertext) {
		return nil, errors.New("rsacrypto: data length error")
	}
	m := new(big.Int).SetBytes(ciphertext)
	if m.Cmp(pub.N) > 0 {
		return nil, errors.New("rsacrypto: data is too large")
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, errors.New("rsacrypto: data broken")
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, errors.New("rsacrypto: decryption error")
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

// DecodePrivKeyPem decode private key in pem format
func DecodePrivKeyPem(privateKey string) (priv *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		err = errors.New("rsacrypto: private key format error")
		return
	}
	priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		err = errors.New("rsacrypto: failed to decode private key: " + err.Error())
		return
	}
	return priv, nil
}

// DecodePubKeyPem decode public key in pem format
func DecodePubKeyPem(publicKey string) (pub *rsa.PublicKey, err error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || (block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY") {
		err = errors.New("rsacrypto: public key format error")
		return
	}
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		err = errors.New("rsacrypto: failed to decode private key: " + err.Error())
		return
	}
	pub, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		err = errors.New("rsacrypto: failed to cast public key to rsa public key")
		return
	}
	return pub, nil
}

// DecodePrivKeyModulus decode private key from modulus
// The modulus parameter is Key modulus, pubExponent parameter is public exponent, privExponent parameter is private exponent
// The parameters are hexadecimal strings
func DecodePrivKeyModulus(modulus, pubExponent, privExponent string) (*rsa.PrivateKey, error) {
	if !isHex(modulus) || !isHex(pubExponent) || !isHex(privExponent) {
		return nil, errors.New("rsacrypto: parameter must be hexadecimal strings")
	}
	n := new(big.Int)
	n.SetString(modulus, 16)
	e := new(big.Int)
	e.SetString(pubExponent, 16)
	d := new(big.Int)
	d.SetString(privExponent, 16)
	privateKey := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: n,
			E: int(e.Int64()),
		},
		D: d,
	}
	return privateKey, nil
}

// DecodePubKeyModulus decode private key from modulus
// The modulus parameter is Key modulus, pubExponent parameter is public exponent
// The parameters are hexadecimal strings
func DecodePubKeyModulus(modulus, pubExponent string) (*rsa.PublicKey, error) {
	if !isHex(modulus) || !isHex(pubExponent) {
		return nil, errors.New("rsacrypto: parameter must be hexadecimal strings")
	}
	n := new(big.Int)
	n.SetString(modulus, 16)
	e := new(big.Int)
	e.SetString(pubExponent, 16)
	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}
	return publicKey, nil
}

func isHex(s string) bool {
	re := regexp.MustCompile(`^(?i)[0-9a-f]+$`)
	return re.MatchString(s)
}

func checkPub(pub *rsa.PublicKey) error {
	if pub.N == nil {
		return errors.New("crypto/rsa: missing public modulus")
	}
	if pub.E < 2 {
		return errors.New("crypto/rsa: public exponent too small")
	}
	if pub.E > 1<<31-1 {
		return errors.New("crypto/rsa: public exponent too large")
	}
	return nil
}

func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)

func decrypt(random io.Reader, priv *rsa.PrivateKey, c *big.Int) (m *big.Int, err error) {
	if c.Cmp(priv.N) > 0 {
		err = errors.New("rsacrypto: decryption error")
		return
	}
	var ir *big.Int
	if random != nil {
		var r *big.Int

		for {
			r, err = rand.Int(random, priv.N)
			if err != nil {
				return
			}
			if r.Cmp(bigZero) == 0 {
				r = bigOne
			}
			var ok bool
			ir, ok = modInverse(r, priv.N)
			if ok {
				break
			}
		}
		bigE := big.NewInt(int64(priv.E))
		rpowe := new(big.Int).Exp(r, bigE, priv.N)
		cCopy := new(big.Int).Set(c)
		cCopy.Mul(cCopy, rpowe)
		cCopy.Mod(cCopy, priv.N)
		c = cCopy
	}
	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}
	if ir != nil {
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}

	return
}

func copyWithLeftPad(dest, src []byte) {
	numPaddingBytes := len(dest) - len(src)
	for i := 0; i < numPaddingBytes; i++ {
		dest[i] = 0
	}
	copy(dest[numPaddingBytes:], src)
}

func modInverse(a, n *big.Int) (ia *big.Int, ok bool) {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)
	g.GCD(x, y, a, n)
	if g.Cmp(bigOne) != 0 {
		return
	}
	if x.Cmp(bigOne) < 0 {
		x.Add(x, n)
	}
	return x, true
}
