package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/encryption/des"
)

func main() {

	str, _ := des.DesEncrypt("sunmi", "sunmi388", "12345678")
	fmt.Println(str)
	str2, _ := des.DesDecrypt(str, "sunmi388", "12345678")
	fmt.Println(str2)

	str, _ = des.DesEncryptECB("sunmi", "sunmi388")
	fmt.Println(str)
	str2, _ = des.DesDecryptECB(str, "sunmi388")
	fmt.Println(str2)

	str, _ = des.TripleDesEncrypt("sunmi", "sunmi388sunmi388sunmi388", "12345678")
	fmt.Println(str)
	str2, _ = des.TripleDesDecrypt(str, "sunmi388sunmi388sunmi388", "12345678")
	fmt.Println(str2)

}
