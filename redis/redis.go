package redis

import (
	"github.com/sunmi-OS/gocore/utils"
	"github.com/sunmi-OS/gocore/viper"
	"gopkg.in/redis.v5"
	"strings"
	"sync"
)

var RedisList sync.Map

func GetRedisOptions(db string) {
	host := viper.GetEnvConfig("redisServer.host")
	port := viper.GetEnvConfig("redisServer.port")
	auth := viper.GetEnvConfig("redisServer.auth")
	encryption := viper.GetEnvConfigInt("redisServer.encryption")
	dbIndex := viper.GetEnvConfigCastInt("redisDB." + db)
	if encryption == 1 {
		auth = utils.GetMD5(auth)
	}
	addr := host + port
	if !strings.Contains(addr, ":") {
		addr = host + ":" + port
	}
	options := redis.Options{Addr: addr, Password: auth, DB: dbIndex}
	client := redis.NewClient(&options)
	client.Ping().Result()

	RedisList.Store(db, client)
}

func GetRedisDB(db string) *redis.Client {
	if v, ok := RedisList.Load(db); ok {
		return v.(*redis.Client)
	}

	return nil
}
