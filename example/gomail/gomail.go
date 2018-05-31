package main

import (
	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/gomail"
)

func main() {

	viper.NewConfig("config", "conf")

	gomail.SendEmail("wenzhenxi@vip.qq.com", "service@sunmi.com", "service", "title", "msg")
}
