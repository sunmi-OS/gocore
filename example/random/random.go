package main

import (
	"gocore/utils"
	"fmt"
)

func main(){

	fmt.Println("RandInt",utils.RandInt(13,233))

	fmt.Println("RandInt64",utils.RandInt64(13,233))

	fmt.Println("GetRandomString",utils.GetRandomString(122))

	fmt.Println("GetRandomNumeral",utils.GetRandomNumeral(133))
}
