package main

import (
	"fmt"
	"github.com/sunmi-OS/gocore/redis"
	"github.com/sunmi-OS/gocore/viper"
)

func main() {
	viper.NewConfig("config", "conf")

	redis.GetRedisOptions("email_check")
	redis.GetRedisDB("email_check").Set("test", "sunmi")
	fmt.Println(redis.GetRedisDB("encryption").Get("test").String())
}
