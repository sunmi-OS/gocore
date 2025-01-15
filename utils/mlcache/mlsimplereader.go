package mlcache

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"gomod.sunmi.com/gomoddepend/appstore-common/pkg/maths"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/sunmi-OS/gocore/v2/api/ecode"
	"github.com/sunmi-OS/gocore/v2/glog"
	"gorm.io/gorm"
)

var (
	ErrInNotString = ecode.NewV2(-1, "input is not string")
)

const DefaultRetryTime = 3

type SimpleReader struct {
	ml        MLCache
	simpleOpt SimpleOpt
}

type SimpleOpt struct {
	Opt            Opt
	NotFoundFunc   func(error) bool
	CacheKeyPrefix string
	Retry          int // 默认3次重试。1代表最多请求1次，2代表最多请求3次，以此类推。一般不需要调整，若<=0则会被调整为3。
}

// SimpleCallback 读mysql/mongodb等数据库 or 其他io
type SimpleCallback func(ctx context.Context, key string) (val interface{}, err error)

func GetCacheKey(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", prefix, key)
}

// NewSimpleReader 需要传入指针类型，如&[]*struct{}
// 先读redis缓存，如果没有再读取callback，然后写入redis缓存
// redis key存储格式是 cacheKeyPrefix:{key}
func NewSimpleReader(client *redis.Client, callback SimpleCallback, simpleOpt ...SimpleOpt) *SimpleReader {
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
			cacheValue, err := client.Get(ctx, cacheKey).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					return nil, false, nil
				}
				glog.ErrorC(ctx, "failed to get cache value key:%s err %s", key, err)
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
			err := client.Set(ctx, cacheKey, value, inTTL).Err()
			if err != nil {
				glog.ErrorC(ctx, "SetCacheHandler failed key:%v cacheKey:%v value:%s err:%v", key, cacheKey, value, err)
			}
			return err
		},
		CleanCacheHandler: func(ctx context.Context, key string) error {
			cacheKey := GetCacheKey(sOpt.CacheKeyPrefix, key)
			err := client.Del(ctx, cacheKey).Err()
			if err != nil {
				glog.ErrorC(ctx, "CleanCacheHandler failed key:%s cacheKey:%s err %s", key, cacheKey, err)
			}
			return err
		},
		Decoder: func(input interface{}, result interface{}) error {
			if value, ok := input.(string); ok {
				return sonic.UnmarshalString(value, result)
			} else {
				return ErrInNotString
			}
		},
		Encoder: func(input interface{}) (interface{}, error) {
			return sonic.MarshalString(input)
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
	return &SimpleReader{
		simpleOpt: sOpt,
		ml:        ml,
	}
}

func (m *SimpleReader) Get(ctx context.Context, key string, value interface{}, opt ...Opt) (CacheStatus, error) {
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

func (m *SimpleReader) Delete(ctx context.Context, key string) error {
	return m.ml.Clean(ctx, key)
}

func GetPointer(value interface{}) interface{} {
	pointerType := reflect.PtrTo(reflect.TypeOf(value))
	pointer := reflect.New(pointerType.Elem())
	pointer.Elem().Set(reflect.ValueOf(value))
	return pointer.Interface()
}

func GetValue(pointer interface{}) interface{} {
	value := reflect.ValueOf(pointer).Elem()
	return value.Interface()
}
