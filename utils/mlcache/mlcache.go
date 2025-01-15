package mlcache

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/rogpeppe/go-internal/cache"
	"github.com/sunmi-OS/gocore/v2/api/ecode"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"
	"golang.org/x/sync/singleflight"
)

/*   requirement collection
redis cache with distributed ttl value to avoid snow crash
[done]use gcore retry
*/

var (
	ErrNotFound            = ecode.NewV2(-1, "not found")
	ErrSetHandleNotFound   = ecode.NewV2(-1, "set handler not found")
	ErrCleanHanderNotFound = ecode.NewV2(-1, "clean handler not found")
	ErrInputNotString      = ecode.NewV2(-1, "input is not string")
)

const (
	LevelL1 = "L1"
	LevelL2 = "L2"
	LevelL3 = "L3"
)

type CacheStatus struct {
	// the ttl of key life cycle
	// ttl time.Duration

	// if found the key
	Found bool

	// if the val staled
	// Stale bool

	// CacheFlag: L1/L2/L3
	CacheFlag string
}

type Opt struct {
	// the ttl of key life cycle
	TTL time.Duration

	Timeout time.Duration
}

type MLCache interface {
	// TODO: support context
	Set(ctx context.Context, key string, val interface{}, opt ...Opt) error
	SetL2(ctx context.Context, key string, val interface{}, opt ...Opt) error
	Get(ctx context.Context, key string, val interface{}, opt ...Opt) (CacheStatus, error)
	Clean(ctx context.Context, key string) error
}

type GetCacheHandler func(ctx context.Context, key string) (val interface{}, found bool, err error)

type SetCacheHandler func(ctx context.Context, key string, val interface{}, ttl *time.Duration) error

type CleanCacheHandler func(ctx context.Context, key string) error

type Encoder func(val interface{}) (str interface{}, err error)

type DecoderType interface {
	~int32 | string | []byte
}

type Decoder func(str interface{}, val interface{}) error

type LC struct {
	Name      string
	LevelName string
	// retry times when get/set handler failed
	Retry             int
	GetCacheHandler   GetCacheHandler
	SetCacheHandler   SetCacheHandler
	CleanCacheHandler CleanCacheHandler
	// encoder and decoder
	Encoder Encoder
	Decoder Decoder
}

func (lc *LC) Get(ctx context.Context, key string, pointer interface{}) (cs CacheStatus, err error) {
	if lc.GetCacheHandler == nil {
		cs.Found = false
		return
	}
	// TODO: add once setup to avoid re-connect function on first moment
	if lc.Decoder == nil {
		lc.Decoder = DefaultDecoder
	}

	var found bool
	var val interface{}
	retry := lc.Retry

	// TODO[Done]: add timeout; do with retry from gocore
	retryFunc := func() error {
		val, found, err = lc.GetCacheHandler(ctx, key)
		return err
	}
	// todo 某些请求，不需要报错，直接返回
	err = utils.Retry(retryFunc, retry, time.Millisecond*200)

	cs.Found = found
	if cs.Found {
		cs.CacheFlag = lc.LevelName
	}
	if err != nil || !found {
		return cs, err
	}

	err = lc.Decoder(val, pointer)
	return cs, err
}

func (lc *LC) Set(ctx context.Context, key string, val interface{}, ttl *time.Duration) (err error) {
	inTTL := time.Duration(0)
	if ttl != nil {
		inTTL = *ttl
	}
	if lc.SetCacheHandler == nil {
		glog.WarnC(ctx, "%v lc.SetCacheHandler: set handler not found", lc.Name)
		return ErrSetHandleNotFound
	}
	if lc.Encoder == nil {
		lc.Encoder = DefaultEncoder
	}

	inputVal, err := lc.Encoder(val)
	if err != nil {
		glog.WarnC(ctx, "%v lc.SetCacheHandler: val:%+v, err:%+v", lc.Name, inputVal, err)
		return err
	}
	return lc.SetCacheHandler(ctx, key, inputVal, &inTTL)
}

func (lc *LC) Clean(ctx context.Context, key string) (err error) {
	if lc.CleanCacheHandler == nil {
		glog.WarnC(ctx, "%v lc.CleanCacheHandler: clean handler not found", lc.Name)
		return ErrCleanHanderNotFound
	}
	return lc.CleanCacheHandler(ctx, key)

}

type mLCache struct {
	// L1 cache ---> go cache left empty for now
	L1 *cache.Cache

	// L2 cache ---> redis cache
	L2 *LC

	// L3 cache ---> mysql
	L3 *LC

	Opt Opt
	// lock cache key
	//Lock *KeyLock
	//
	//// global mutex
	//Mu *sync.Mutex

	SingleHandle singleflight.Group
}

func New(
	l2, l3 *LC,
) MLCache {
	if l2 != nil {
		l2.LevelName = LevelL2
	}
	if l3 != nil {
		l3.LevelName = LevelL3
	}
	return &mLCache{
		L1: nil,
		//Lock:  newKeyLock(0),
		L2: l2,
		L3: l3,
		//Mu:    &sync.Mutex{},
		SingleHandle: singleflight.Group{},
	}
}

// TODO
// 1. set l3 cache
// 2. set l2 cache
// 3. set l1 cache
// func (mlc *mLCache) Set(key string, val interface{}, opt Opt) { }

func (mlc *mLCache) Get(ctx context.Context, key string, pointer interface{}, opt ...Opt) (cs CacheStatus, err error) {
	tOpt := DefaultOpt(mlc.Opt, opt...)
	// err can not be nil

	// L1 cache left empty
	//val, cs, err = mlc.GetFromL1Cache(key, ctx)
	//if err != nil {
	//	return
	//}
	// hit l1 cache
	//if cs.Found && !cs.Stale {
	//	cs.CacheFlag = "L1"
	//	return
	//}
	// missing L1 cache
	// first: fetch from L1 cache

	// no L2 cache, should not let val/cs/err be covered
	var val interface{}
	val, err, _ = mlc.SingleHandle.Do(key, func() (result interface{}, err error) {
		// first: try get cache
		if mlc.L2 == nil {
			goto L3
		}
		cs, err = mlc.L2.Get(ctx, key, pointer)
		if err != nil {
			// glog.WarnC(ctx, "call L2 %v Get, key:%s, err:%+v", mlc.L2.Name, key, err)
			goto L3
		}
		// hit l2 cache
		if cs.Found {
			return pointer, nil
		}

	L3:
		if mlc.L3 == nil {
			glog.WarnC(ctx, "L3 nil Get, key:%s skip", key)
			return nil, nil
		}
		// second: fetch from L2 cache and set to L1 cache
		cs, err = mlc.L3.Get(ctx, key, pointer)
		if err != nil {
			glog.WarnC(ctx, "call L3 %v Get, key:%s, err:%+v", mlc.L3.Name, key, err)
			return nil, err
		}
		if !cs.Found {
			return nil, nil
		}
		// try set L2 cache when L3 value get
		if mlc.L2 != nil {
			err = mlc.L2.Set(ctx, key, pointer, &tOpt.TTL)
			if err != nil {
				glog.WarnC(ctx, "call L2 %v Set, key:%s, err:%+v", mlc.L2.Name, key, err)
				return nil, err
			} else {
				return pointer, nil
			}
		}
		return nil, nil
	})
	if val != nil && val != pointer {
		err = CopyInterface(val, pointer)
		if err == nil {
			cs.Found = true
		}
	}
	return cs, err
}

func (mlc *mLCache) Set(ctx context.Context, key string, val interface{}, opt ...Opt) (err error) {
	tOpt := DefaultOpt(mlc.Opt, opt...)
	if mlc.L3 != nil {
		err = mlc.L3.Set(ctx, key, val, &tOpt.TTL)
		if err != nil {
			return
		}
		// clean L2 cache
		if mlc.L2 != nil {
			err := mlc.L2.Clean(ctx, key)
			if err != nil {
				glog.WarnC(ctx, "call L2.Clean, key:%s, err:%+v", key, err)
				return err
			}
		}
		return
	}
	if mlc.L3 == nil {
		err = mlc.L2.Set(ctx, key, val, &tOpt.TTL)
		return
	}

	return
}

func (mlc *mLCache) SetL2(ctx context.Context, key string, val interface{}, opt ...Opt) (err error) {
	tOpt := DefaultOpt(mlc.Opt, opt...)
	if mlc.L2 != nil {
		err = mlc.L2.Set(ctx, key, val, &tOpt.TTL)
		if err != nil {
			glog.WarnC(ctx, "call L2.Set, key:%s, err:%+v", key, err)
			return
		}
	}
	return
}

func (mlc *mLCache) Clean(ctx context.Context, key string) error {
	if mlc.L2 != nil {
		err := mlc.L2.Clean(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func DefaultEncoder(input interface{}) (result interface{}, err error) {
	return sonic.MarshalString(input)
}

func DefaultDecoder(input interface{}, result interface{}) (err error) {
	if value, ok := input.(string); ok {
		return sonic.UnmarshalString(value, result)
	} else {
		return ErrInputNotString
	}
}

func DefaultOpt(origin Opt, in ...Opt) Opt {
	op := Opt{}
	if len(in) != 0 {
		op = in[0]
	} else {
		op = origin
	}
	return op
}

func DefaultTTL(ttl *time.Duration) time.Duration {
	if ttl != nil {
		return *ttl
	}
	return time.Duration(0)
}
