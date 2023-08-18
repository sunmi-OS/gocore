package redis

import (
	"context"
	"crypto/tls"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils/closes"
	"github.com/sunmi-OS/gocore/v2/utils/hash"
)

var Map sync.Map
var closeOnce sync.Once

// NewRedis new Redis Client
func NewRedis(dbName string) {
	rc, err := newRedis(dbName)
	if err != nil {
		panic(err)
	}
	Map.Store(dbName, rc)
	closeOnce.Do(func() {
		closes.AddShutdown(closes.ModuleClose{
			Name:     "Redis Close",
			Priority: closes.RedisPriority,
			Func:     Close,
		})
	})
}

func newRedis(db string) (rc *redis.Client, err error) {
	redisName, dbName := splitDbName(db)

	host := viper.GetEnvConfig(redisName + ".host").String()
	port := viper.GetEnvConfig(redisName + ".port").String()
	auth := viper.GetEnvConfig(redisName + ".auth").String()
	encryption := viper.GetEnvConfig(redisName + ".encryption").Int64()
	dbIndex := viper.GetEnvConfig(redisName + ".redisDB." + dbName).Int()
	insecureSkipVerify := viper.GetEnvConfig(redisName + ".insecureSkipVerify").Bool()
	if redisName == "redisServer" {
		dbIndex = viper.GetEnvConfig("redisDB." + dbName).Int()
	}
	if encryption == 1 {
		auth = hash.MD5(auth)
	}
	addr := host + port
	if !strings.Contains(addr, ":") {
		addr = host + ":" + port
	}
	ops := &redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       dbIndex,
	}
	if insecureSkipVerify {
		ops.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	rc = redis.NewClient(ops)
	if err := rc.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rc, nil
}

// NewOrUpdateRedis 新建或更新redis客户端
func NewOrUpdateRedis(dbName string) error {
	rc, err := newRedis(dbName)
	if err != nil {
		return err
	}

	v, _ := Map.Load(dbName)
	Map.Store(dbName, rc)

	if v != nil {
		err := v.(*redis.Client).Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetRedis 获取 RedisClient
func GetRedis(dbName string) *redis.Client {
	if v, ok := Map.Load(dbName); ok {
		return v.(*redis.Client)
	}
	return nil
}

func Close() {
	Map.Range(func(dbName, rc interface{}) bool {
		glog.WarnF("close db %s", dbName)
		Map.Delete(dbName)
		err := rc.(*redis.Client).Close()
		return err == nil
	})
}

func splitDbName(db string) (redisName, dbName string) {
	kv := strings.Split(db, ".")
	if len(kv) == 2 {
		return kv[0], kv[1]
	}
	if len(kv) == 1 {
		return "redisServer", kv[0]
	}
	panic("redis dbName Mismatch")
}
