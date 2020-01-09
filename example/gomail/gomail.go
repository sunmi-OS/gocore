package main

import (
	"github.com/sunmi-OS/gocore/gomail"
	"github.com/sunmi-OS/gocore/viper"
)

func main() {

	viper.NewConfig("config", "conf")

	gomail.SendEmail("wenzhenxi@vip.qq.com", "service@sunmi.com", "service", "title", "msg")
}
