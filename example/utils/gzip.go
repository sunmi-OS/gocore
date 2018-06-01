package main

import (
	"gocore/utils"
	"fmt"
)

func main(){
	fmt.Println(utils.GzipEncode("dsxdjdhskfjkdsfhsdjlaal"))
	var m = utils.GzipEncode("dsxdjdhskfjkdsfhsdjlaal")
	fmt.Println(utils.GzipDecode(m))
}
