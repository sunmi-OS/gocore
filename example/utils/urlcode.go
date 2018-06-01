package main

import (
	"gocore/utils"
	"fmt"
)

func main(){

	var urls string
	urls = "https://www.sunmi.com/"
	e,err:= utils.UrlEncode(urls)
	if err != nil{
		fmt.Println("UrlEncode failed error",err)
	}

	fmt.Println("UrlEncode",e)

	r,err := utils.UrlDecode(urls)
	if err != nil{
		fmt.Println("UrlDecode failed error",err)
	}
	fmt.Println("UrlDecode",r)
}
