package redis

import (
	"strings"
	"sync"

	viper2 "github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/utils"
	"gopkg.in/redis.v5"
)

var RedisList sync.Map

// Deprecated
// 推荐使用 NewRedis
func GetRedisOptions(db string) {

	client, err := openRedis(db)
	if err != nil {
		panic(err)
	}

	RedisList.Store(db, client)
}

// Deprecated
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

// Deprecated
func GetRedisDB(db string) *redis.Client {
	if v, ok := RedisList.Load(db); ok {
		return v.(*redis.Client)
	}

	return nil
}

func openRedis(db string) (*redis.Client, error) {

	redisName, dbName := dbNameSplit(db)

	host := viper2.GetEnvConfig(redisName + ".host")
	port := viper2.GetEnvConfig(redisName + ".port")
	auth := viper2.GetEnvConfig(redisName + ".auth")
	encryption := viper2.GetEnvConfigInt(redisName + ".encryption")
	dbIndex := viper2.GetEnvConfigCastInt(redisName + ".redisDB." + dbName)
	if redisName == "redisServer" {
		dbIndex = viper2.GetEnvConfigCastInt("redisDB." + dbName)
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
	}
	if len(kv) == 1 {
		return "redisServer", kv[0]
	}
	panic("redis dbName Mismatch")
}
