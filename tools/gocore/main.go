package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	constant "github.com/sunmi-OS/gocore/tools/gocore/constant"
	"github.com/sunmi-OS/gocore/viper"
)

func main() {
	fmt.Printf("%#v\n", create())
}

func create() error {
	viper.NewConfig("config", "config")
	service := viper.C.GetStringMap("service")
	name := viper.C.GetString("service.name")
	err := MkdirIfNotExist(name)
	if err != nil {
		panic(err)
	}
	path := name + "/main.go"
	file, e := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	if e != nil {
		return e
	}
	defer file.Close()
	_, err = io.WriteString(file, constant.MainTemplate)
	if err != nil {
		panic(err)
	} else {
		exec.Command("goimports", "-l", "-w", path).Output()
		log.Println(path + " 已生成...")
	}
	fmt.Printf("%#v\n", service)
	return nil
}
