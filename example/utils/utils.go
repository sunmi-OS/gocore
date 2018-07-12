package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/utils"
)

func main() {

	d := utils.GetDate()
	fmt.Println("GetData", d)

	m := utils.GetRunTime()
	fmt.Println("GetRunTime", m)

	var encryption string
	encryption = "1243sdfds"

	t := utils.GetMD5(encryption)
	fmt.Println("GetMD5", t)
}
