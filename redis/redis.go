package redis

import (
	"strings"
	"sync"

	"github.com/sunmi-OS/gocore/utils"
	"github.com/sunmi-OS/gocore/viper"

	"gopkg.in/redis.v5"
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

	err = v.(*redis.Client).Close()
	if err != nil {
		return err
	}
	return nil

}

func GetRedisDB(db string) *redis.Client {
	if v, ok := RedisList.Load(db); ok {
		return v.(*redis.Client)
	}

	return nil
}

func openRedis(db string) (*redis.Client, error) {

	redisName, dbName := dbNameSplit(db)

	host := viper.GetEnvConfig(redisName + ".host")
	port := viper.GetEnvConfig(redisName + ".port")
	auth := viper.GetEnvConfig(redisName + ".auth")
	encryption := viper.GetEnvConfigInt(redisName + ".encryption")

	if redisName == "redisServer" {
		dbIndex := viper.GetEnvConfigCastInt("redisDB." + dbName)
	} else {
		dbIndex := viper.GetEnvConfigCastInt(redisName + ".redisDB." + dbName)
	}

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

func dbNameSplit(db string) (redisName, dbName string) {

	kv := strings.Split(db, ".")
	if len(kv) == 2 {
		return kv[0], kv[1]
	} else if len(kv) == 1 {
		return "redisServer", kv[0]
	} else {
		panic("redis dbName Mismatch")
	}
}
