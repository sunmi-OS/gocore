package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/httplib"
)

func main() {

	b := httplib.Post("https://baidu.com/")
	b.Param("username", "astaxie")
	b.Param("password", "123456")
	b.PostFile("uploadfile1", "httplib.pdf")
	b.PostFile("uploadfile2", "httplib.txt")
	str, err := b.String()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
}
