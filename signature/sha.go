package signature

import (
    "crypto/sha256"
    "encoding/hex"
)

func Sha256(data string, salt string) string {
    hash := sha256.New()
    hash.Write([]byte(salt + data))
    hashString := hash.Sum(nil)
    return hex.EncodeToString(hashString)
}
