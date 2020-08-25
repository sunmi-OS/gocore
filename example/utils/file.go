package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/utils"
)

func main() {
	fmt.Println("GetPath:", utils.GetPath())
	fmt.Println(utils.IsDirExists("/tmp/go-build803419530/command-line-arguments/_obj/exe"))
	fmt.Println(utils.MkdirFile("test"))
}
