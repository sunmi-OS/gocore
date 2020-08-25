package main

import (
	"fmt"
	"time"

	"github.com/sunmi-OS/gocore/aliyunlog"
	"github.com/sunmi-OS/gocore/viper"
)

func main() {
	viper.C.SetDefault("log.Project", "xxxx")
	viper.C.SetDefault("log.Endpoint", "cn-hangzhou-intranet.log.aliyuncs.com")
	viper.C.SetDefault("log.AccessKey", "xxxxx")
	viper.C.SetDefault("log.SecretKey", "xxxxx")
	viper.C.SetDefault("log.LogStore", "xxxx")

	aliyunlog.InitLog("log", "logStore")

	for i := 0; i < 10; i++ {
		aliyunlog.Info("test", map[string]string{"content": "test", "content2": fmt.Sprintf("%v", i)})
	}

	time.Sleep(1 * time.Second)
}
