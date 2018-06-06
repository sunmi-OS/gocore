package main

import (
	"github.com/sunmi-OS/gocore/utils"
	"fmt"
)

func main(){
	var secret,params string
	secret = "123"
	params = "abdsjfhdshfksdhf"
	m,err := utils.GetParamMD5Sign(secret,params)
	if err !=nil{
		fmt.Println("Error:",err)
	}
	fmt.Println("GetParamMD5Sign",m)

	var maintain string

	maintain = "dfssjfdsdfghjdsfgdsj"
	n,err := utils.GetSHA(maintain)
	if err != nil{
		fmt.Println("GetSHA failed error",err)
	}
	fmt.Println("GetSHA",n)

	l,err:= utils.GetParamHmacSHA256Sign(secret,params)
	if err != nil{
		fmt.Println("GetParamHmacSHA256Sign failed err",err)
	}

	fmt.Println("GetParamHmacSHA256Sign",l)

	p,err := utils.GetParamHmacSHA512Sign(secret,params)
	if err != nil{
		fmt.Println("GetParamHmacSHA512Sign failed error",err)
	}
	fmt.Println("GetParamHmacSHA512Sign",p)

	u,err := utils.GetParamHmacSHA1Sign(secret,params)
	if err != nil{
		fmt.Println("GetParamHmacSHA1Sign failed error",err)
	}

	fmt.Println("GetParamHmacSHA1Sign",u)

	c,err := utils.GetParamHmacMD5Sign(secret,params)
	if err != nil{
		fmt.Println("GetParamHmacMD5Sign failed error",err)
	}

	fmt.Println("GetParamHmacMD5Sign",c)


	d,err := utils.GetParamHmacSha384Sign(secret,params)
	if err != nil{
		fmt.Println("GetParamHmacSha384Sign failed error",err)
	}

	fmt.Println("GetParamHmacSha384Sign",d)

	f,err := utils.GetParamHmacSHA256Base64Sign(secret,params)
	if err != nil{
		fmt.Println("GetParamHmacSHA256Base64Sign failed error",err)
	}

	fmt.Println("GetParamHmacSHA256Base64Sign",f)

	var hmac_key,hmac_data string
	hmac_key = "12322334234"
	hmac_data = "sjhdjsdjfh"
	t:= utils.GetParamHmacSHA512Base64Sign(hmac_key,hmac_data)

	fmt.Println("GetParamHmacSHA512Base64Sign",t)
}
