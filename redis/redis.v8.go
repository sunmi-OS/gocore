package redis

import (
	"context"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/sunmi-OS/gocore/utils"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/xlog"
)

type Client struct {
	redisMaps sync.Map
}

// NewRedis new Redis Client
func NewRedis(dbName string) (c *Client) {
	c = new(Client)
	rc, err := newRedis(dbName)
	if err != nil {
		panic(err)
	}
	c.redisMaps.Store(dbName, rc)
	return c
}

func newRedis(db string) (rc *redis.Client, err error) {
	redisName, dbName := splitDbName(db)

	host := viper.GetEnvConfig(redisName + ".host")
	port := viper.GetEnvConfig(redisName + ".port")
	auth := viper.GetEnvConfig(redisName + ".auth")
	encryption := viper.GetEnvConfigInt(redisName + ".encryption")
	dbIndex := viper.GetEnvConfigCastInt(redisName + ".redisDB." + dbName)
	if redisName == "redisServer" {
		dbIndex = viper.GetEnvConfigCastInt("redisDB." + dbName)
	}
	if encryption == 1 {
		auth = utils.GetMD5(auth)
	}
	addr := host + port
	if !strings.Contains(addr, ":") {
		addr = host + ":" + port
	}

	rc = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       dbIndex,
	})
	if err := rc.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rc, nil
}

// NewOrUpdateRedis 新建或更新redis客户端
func (c *Client) NewOrUpdateRedis(dbName string) error {
	rc, err := newRedis(dbName)
	if err != nil {
		return err
	}

	v, _ := c.redisMaps.Load(dbName)
	c.redisMaps.Delete(dbName)
	c.redisMaps.Store(dbName, rc)

	if v != nil {
		v.(*redis.Client).Close()
	}
	return nil
}

// GetRedis 获取 RedisClient
func (c *Client) GetRedis(dbName string) *redis.Client {
	if v, ok := c.redisMaps.Load(dbName); ok {
		return v.(*redis.Client)
	}
	return nil
}

func (c *Client) Close() {
	c.redisMaps.Range(func(dbName, rc interface{}) bool {
		xlog.Warnf("close db %s", dbName)
		c.redisMaps.Delete(dbName)
		rc.(*redis.Client).Close()
		return true
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
