# gocore-mlcache

Fast and automated layered caching for Golang.

This library can be manipulated as a key/value store caching golang types , combining the power of the fastcache and redis, which results in an extremely performant and flexible caching solution.

Features:

- Caching with TTLs.
- Use singleflight to prevent dog-pile effects to your
  database/backend on cache misses.
- Support prometheus metrics

Illustration of the various caching levels built into this library:


```
┌─────────────────────────────────────────────────┐
│ golang                                          │
│       ┌───────────────────────────────────────┐ │
│       │                                       │ │
│ L1    │                                       │ │
│       │             memory cache              │ │
│       └───────────────────────────────────────┘ │
│             │             │             │       │
│             ▼             ▼             ▼       │
│       ┌───────────────────────────────────────┐ │
│       │                                       │ │
│ L2    │              redis cache              │ │
│       │                                       │ │
│       └───────────────────────────────────────┘ │
│                           │ singleflight        │
│                           ▼                     │
│                  ┌──────────────────┐           │
│                  │     callback     │           │
│                  └────────┬─────────┘           │
└───────────────────────────┼─────────────────────┘
                            │
  L3                        │   I/O fetch
                            ▼

                   Database, API, DNS, Disk, any I/O...
```


The cache level hierarchy is:
- **L1**: Least-Recently-Used golang memory cache using fastcache
- **L2**: Redis cache. This level is only accessed if L1 was a miss, 
  and prevents workers from requesting the L3 cache.
- **L3**: A custom function that will only be run by a single goroutine
  to avoid the dog-pile effect on your database/backend
  (via singleflight). Values fetched via L3 will be set to the L2 cache
  for other goroutine to retrieve.

## Usage
### NewSimpleReader
```golang
package lcache

import (
  "context"
  "fmt"

  "xxx/dal/cache"
  "xxx/dal/devicedb"
  "xxx/dal/devicedb/model"
  "github.com/sunmi-OS/gocore/v2/pkg/mlcache"
)

const DeviceInfoCachePrefix = "di"

func main() {
  var deviceInfoClient *redis.Client
  // init deviceInfoClient
  deviceInfoML := mlcache.NewSimpleReader(deviceInfoClient, func(ctx context.Context, key string) (interface{}, error) {
    return devicedb.GetDeviceInfoByID(ctx, key), nil
  }, mlcache.SimpleOpt{
    NotFoundFunc:   IsErrNotFound, // 除了gorm.ErrRecordNotFound外，若有其他错误需要跳过重试，定义在这里
    CacheKeyPrefix: DeviceInfoCachePrefix,       // 如果定义了CacheKeyPrefix，那么redis存储时key就会加上前缀，变成 CacheKeyPrefix:key
  })
  info := &model.IotDevice{}
  ctx := context.Background()
  status, err := deviceInfoML.Get(ctx, "device_id", info)
  fmt.Printf("%+v\n", status)
  fmt.Println(err)
  fmt.Printf("%+v\n", info)
}

```

### NewMemoryDBReader

```
import(
  gomod.sunmi.com/gomoddepend/appstore-common/pkg/mlcache"
  ristretto_store "github.com/eko/gocache/store/ristretto/v4"
)

var deviceInfoML *mlcache.MemoryDBReader

func GetDeviceInfo(ctx context.Context, id int64) (*model.DeviceInfo, error) {
	data := &model.DeviceInfo{}
	_, err := deviceInfoML.Get(ctx, cast.ToString(vid), data)
	return data, err
}

func main() {
  deviceInfoSmallRCache, err := ristretto.NewCache(&ristretto.Config{
      NumCounters: 2e5,       // 需要是持久化数量的10倍，这里设置为20w
      MaxCost:     100 << 20, // 存储字节数
      BufferItems: 64,
  })
  if err != nil {
      panic(err)
  }
  deviceInfoSmallRStore = ristretto_store.NewRistretto(deviceInfoSmallRCache)
  deviceInfoML = mlcache.NewMemoryDBReader(deviceInfoSmallRStore, func(ctx context.Context, key string) (val interface{}, err error) {
      return GetDeviceInfoByHttp(ctx, cast.ToInt64(key))
  }, mlcache.SimpleOpt{
      CacheKeyPrefix: vidInfoPrefix,
      Opt:            mlcache.Opt{TTL: deviceInfoExpireTTL},
  })

  // 数据预热，提前加载
  datas, err := GetManyDeviceInfos(ctx, deviceIds)
  if err != nil {
      return
  }
  for _, data := range datas {
      _ = deviceInfoML.Set(ctx, cast.ToString(data.Vid), data)
  }
}

```


## TODO
- [ ] 数据统计：缓存命中率、缓存回源率、读写速度、读写qps
- [x] 支持内存缓存 