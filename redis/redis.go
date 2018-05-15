package redis

import (
	"gopkg.in/redis.v5"

	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/utils"
)

var RedisList map[string]*redis.Client

func init() {
	RedisList = make(map[string]*redis.Client, 10)
}

func GetRedisOptions(db string) {
	host := viper.C.GetString("redisServer.host")
	post := viper.C.GetString("redisServer.port")
	auth := viper.C.GetString("redisServer.auth")
	//prefix := viper.C.GetString("redisServer.prefix")
	encryption := viper.C.GetInt("redisServer.encryption")
	dbIndex := viper.C.GetInt("redisDB." + db)
	if encryption == 1 {
		auth = utils.GetMD5(auth)
	}
	options := redis.Options{Addr: host + post, Password: auth, DB: dbIndex,}
	client := redis.NewClient(&options)
	client.Ping().Result()
	RedisList[db] = client
}

func GetRedisDB(db string) *redis.Client {
	if _, ok := RedisList[db]; ok {
		//存在
		return RedisList[db]
	}
	return nil
}
