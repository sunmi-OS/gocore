package main

import (
	"fmt"
	
	"github.com/sunmi-OS/gocore/viper"
)

func main() {
	// 指定配置文件所在的目录和文件名称
	viper.NewConfig("config", "conf")
	// 打印读取的配置
	fmt.Println("port : ", viper.C.Get("system.port"))
	fmt.Println("ENV RUN_TIME : ", viper.GetEnvConfig("run.time"))
}
