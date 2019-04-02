package signature

import (
    "crypto/md5"
    "encoding/hex"
)

func Md5(data string) string {
    m := md5.New()
    m.Write([]byte(data))
    sign := m.Sum(nil)
    return hex.EncodeToString(sign)
}
