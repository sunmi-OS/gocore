package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/utils"
)

func main() {
	fmt.Println(utils.GzipEncode("dsxdjdhskfjkdsfhsdjlaal"))
	var m = utils.GzipEncode("dsxdjdhskfjkdsfhsdjlaal")
	fmt.Println(utils.GzipDecode(m))
}
