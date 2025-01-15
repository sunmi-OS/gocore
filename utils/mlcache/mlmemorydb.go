package mlcache

import (
	"context"
	"errors"
	"time"

	"gomod.sunmi.com/gomoddepend/appstore-common/pkg/maths"

	lib_store "github.com/eko/gocache/lib/v4/store"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
	"github.com/sunmi-OS/gocore/v2/glog"
	"gorm.io/gorm"
)

type MemoryDBReader struct {
	ml        MLCache
	simpleOpt SimpleOpt
}

// NewMemoryDBReader callback中val需要是指针类型，如&[]*struct{}
// 先读redis缓存，如果没有再读取callback，然后写入redis缓存
// redis key存储格式是 cacheKeyPrefix:{key}
func NewMemoryDBReader(client *ristretto_store.RistrettoStore, callback SimpleCallback, simpleOpt ...SimpleOpt) *MemoryDBReader {
	sOpt := SimpleOpt{}
	if len(simpleOpt) != 0 {
		sOpt = simpleOpt[0]
	}
	if sOpt.Retry <= 0 {
		sOpt.Retry = DefaultRetryTime
	}
	sOpt.Opt.TTL = maths.ShakeTime10(sOpt.Opt.TTL)

	l3GetHandler := func(ctx context.Context, key string) (val interface{}, found bool, err error) {
		val, err = callback(ctx, key)
		if err == nil {
			return val, true, nil
		}

		if sOpt.NotFoundFunc != nil && sOpt.NotFoundFunc(err) {
			err = nil
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) { // 兼容默认的mysql读找不到
			err = nil
		}
		return
	}

	ml := New(&LC{
		Retry: sOpt.Retry,
		GetCacheHandler: func(ctx context.Context, key string) (interface{}, bool, error) {
			cacheKey := GetCacheKey(sOpt.CacheKeyPrefix, key)
			cacheValue, err := client.Get(ctx, cacheKey)
			if err != nil {
				if !errors.Is(err, &lib_store.NotFound{}) {
					glog.ErrorC(ctx, "failed to get cache value key:%s err %s", key, err)
				} else {
					return nil, false, nil
				}
				return nil, false, err
			}
			return cacheValue, true, nil
		},
		SetCacheHandler: func(ctx context.Context, key string, value interface{}, ttl *time.Duration) error {
			inTTL := sOpt.Opt.TTL
			if ttl != nil {
				inTTL = *ttl
			}
			cacheKey := GetCacheKey(sOpt.CacheKeyPrefix, key)
			err := client.Set(ctx, cacheKey, value, lib_store.WithExpiration(inTTL))
			if err != nil {
				glog.ErrorC(ctx, "SetCacheHandler failed key:%v cacheKey:%v value:%s err:%v", key, cacheKey, value, err)
			}
			return err
		},
		CleanCacheHandler: func(ctx context.Context, key string) error {
			cacheKey := GetCacheKey(sOpt.CacheKeyPrefix, key)
			err := client.Delete(ctx, cacheKey)
			if err != nil {
				glog.ErrorC(ctx, "CleanCacheHandler failed key:%s cacheKey:%s err %s", key, cacheKey, err)
			}
			return err
		},
		Decoder: func(input interface{}, result interface{}) error {
			return CopyInterface(input, result)
		},
		Encoder: func(input interface{}) (interface{}, error) {
			return input, nil
		},
	}, &LC{
		Retry:           sOpt.Retry,
		GetCacheHandler: l3GetHandler,
		SetCacheHandler: func(ctx context.Context, key string, value interface{}, ttl *time.Duration) error {
			return nil
		},
		Decoder: func(input interface{}, result interface{}) error {
			return CopyInterface(input, result)
		},
	})
	return &MemoryDBReader{
		simpleOpt: sOpt,
		ml:        ml,
	}
}

func (m *MemoryDBReader) Get(ctx context.Context, key string, value interface{}, opt ...Opt) (CacheStatus, error) {
	op := Opt{}
	if len(opt) != 0 {
		op = opt[0]
	} else {
		op = m.simpleOpt.Opt
	}
	status, err := m.ml.Get(ctx, key, value, op)
	if err != nil {
		return status, err
	}
	if !status.Found {
		return status, ErrNotFound
	}
	return status, err
}

func (m *MemoryDBReader) Set(ctx context.Context, key string, value interface{}, opt ...Opt) error {
	op := Opt{}
	if len(opt) != 0 {
		op = opt[0]
	} else {
		op = m.simpleOpt.Opt
	}
	return m.ml.SetL2(ctx, key, value, op)
}

func (m *MemoryDBReader) Delete(ctx context.Context, key string) error {
	return m.ml.Clean(ctx, key)
}
