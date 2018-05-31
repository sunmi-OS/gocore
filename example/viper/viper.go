package main

import (
	"fmt"
	
	"github.com/sunmi-OS/gocore/viper"
)

func main() {
	viper.NewConfig("config", "conf")
	fmt.Println("port : ", viper.C.Get("system.port"))
	fmt.Println("ENV RUN_TIME : ", viper.GetEnvConfig("run.time"))
}
