package main

import (
	"fmt"
	"gocore/example/nacos/config"
	"gocore/viper"
	"time"

	"github.com/sunmi-OS/gocore/nacos"
)

func main() {
	config.InitNacos("dev")

	nacos.SetDataIds("DEFAULT_GROUP", "adb")
	nacos.SetDataIds("pay", "test")


	nacos.SetCallBackFunc("DEFAULT_GROUP", "adb", func(namespace, group, dataId, data string) {

		s := viper.C.GetString("remotemanageDB.dbHost")

		fmt.Println(s)
	})

	nacos.NacosToViper()

	s := viper.C.GetString("remotemanageDB.dbHost")

	fmt.Println(s)

	s = viper.C.GetString("redisDB.remote_control")

	fmt.Println(s)


	time.Sleep(time.Second * 1000)

}
