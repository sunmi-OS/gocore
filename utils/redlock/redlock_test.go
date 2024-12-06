package redlock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

const (
	EsLockKey        = "redlock:es:%s"  // redlock:es:{}
	EsLockExpireTime = time.Second * 60 // 统一默认过期时间
)

func TestRedlock(t *testing.T) {
	var c *redis.Client
	lock := New(c)
	key := fmt.Sprintf(EsLockKey, "id")
	esLock := lock.NewMutex(key, redsync.WithExpiry(EsLockExpireTime))

	// lock with Lock
	err := esLock.Lock()
	if err != nil {
		t.Error(err)
	}
	// do sth
	_, _ = esLock.Unlock()

	// lock with LockContext
	ctx := context.Background()
	if err = esLock.LockContext(ctx); err != nil {
		t.Error(err)
	}
	if _, err = esLock.UnlockContext(ctx); err != nil {
		t.Error(err)
	}
}
