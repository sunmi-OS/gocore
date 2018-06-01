package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/encryption/aes"
)

func main() {

	str, _ := aes.AesEncrypt("sunmi", "sunmiWorkOnesunmiWorkOne")
	fmt.Println(str)
	str2, _ := aes.AesDecrypt(str, "sunmiWorkOnesunmiWorkOne")
	fmt.Println(str2)
}
