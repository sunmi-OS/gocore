package main

import (
	"fmt"

	"gocore/redis"

	"github.com/sunmi-OS/gocore/viper"
)

func main() {
	viper.NewConfig("config", "conf")
	redis.GetRedisOptions("e_invoice")
	redis.GetRedisDB("e_invoice").Set("test", "sunmi", 0)
	fmt.Println(redis.GetRedisDB("e_invoice").Get("test").String())

	redis.GetRedisOptions("OtherRedisServer.e_invoice")
	redis.GetRedisDB("OtherRedisServer.e_invoice").Set("test", "sunmi_other", 0)
	fmt.Println(redis.GetRedisDB("OtherRedisServer.e_invoice").Get("test").String())
	fmt.Println(redis.GetRedisDB("e_invoice").Get("test").String())
}
