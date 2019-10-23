package redis

import (

	"strings"
	"sync"

	"gopkg.in/redis.v5"
	"github.com/sunmi-OS/gocore/utils"
	"github.com/sunmi-OS/gocore/viper"
)

var RedisList sync.Map

func GetRedisOptions(db string) {

	client, err := openRedis(db)
	if err != nil {
		panic(err)
	}

	RedisList.Store(db, client)
}

func UpdateRedis(db string) error {

	v, _ := RedisList.Load(db)

	client, err := openRedis(db)
	if err != nil {
		return err
	}

	RedisList.Delete(db)
	RedisList.Store(db, client)

	err := v.(*redis.Client).Close()
	if err != nil {
		return err
	}

}

func GetRedisDB(db string) *redis.Client {
	if v, ok := RedisList.Load(db); ok {
		return v.(*redis.Client)
	}

	return nil
}

func openRedis(db string) (*redis.Client, error) {
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
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
