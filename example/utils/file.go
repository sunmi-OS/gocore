package main

import (
	"github.com/sunmi-OS/gocore/utils"
	"fmt"
)

func main(){
	fmt.Println("GetPath:%s",utils.GetPath())
	fmt.Println(utils.IsDirExists("/tmp/go-build803419530/command-line-arguments/_obj/exe"))
	fmt.Println(utils.MkdirFile("test"))
}
