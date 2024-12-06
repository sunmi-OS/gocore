package redlock

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type Redlock struct {
	redsync *redsync.Redsync
}

// New 简单封装，仅支持单个redlock
func New(c *redis.Client) *Redlock {
	pool := goredis.NewPool(c)
	rs := redsync.New(pool)
	return &Redlock{
		redsync: rs,
	}
}

func (r *Redlock) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return r.redsync.NewMutex(name, options...)
}
