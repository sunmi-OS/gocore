package main

import (
	"fmt"

	"github.com/spf13/cast"
)

func main() {

	var i64 int64
	i64 = 60

	toString(cast.ToString(i64))

}

func toString(s string) {

	fmt.Println("这是一个字符串:", s)
}
