package redis

import (
	"sync"

	"gopkg.in/redis.v5"

	"github.com/sunmi-OS/gocore/utils"
	"github.com/sunmi-OS/gocore/viper"
)

var RedisList sync.Map

func GetRedisOptions(db string) {
	host := viper.C.GetString("redisServer.host")
	post := viper.C.GetString("redisServer.port")
	auth := viper.C.GetString("redisServer.auth")
	encryption := viper.C.GetInt("redisServer.encryption")
	dbIndex := viper.C.GetInt("redisDB." + db)
	if encryption == 1 {
		auth = utils.GetMD5(auth)
	}
	options := redis.Options{Addr: host + ":" + post, Password: auth, DB: dbIndex}
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
