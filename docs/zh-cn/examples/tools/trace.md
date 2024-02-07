# 链路追踪

使用脚手架生成的项目默认集成了 Opentelemetry 链路追踪，propagators 默认支持是 Zipkin 的 B3。

使用方式如下：

## HTTP api 

修改 api 层的 handler 方法：

```go
package account

import (
    "xxx/app/biz"
  
	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/v2/api"
	tracing "github.com/sunmi-OS/gocore/v2/lib/tracing/client/otel"
）


func (m *User) Login(g *gin.Context) {
	ctx := api.NewContext(g)
  
  // 链路追踪
  spanCtx, span := tracing.StartSpan(ctx.Request.Context())
  defer span.End()
  
  // 传递 spanCtx
  token, err := biz.UserHandler.Login(spanCtx, username, password)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Success(token)
 }
```

修改 biz 层对应的方法：

```go
package biz

import (
	"context"
)

// 如果需要追踪 biz 层链路，需要传递 context
func (m *User) Login(ctx context.Context, username, password string) (token string, err error) {{
	//...
	
}
```

## Grpc

- gprc 客户端使用方法

```go
import (
	tracing "github.com/sunmi-OS/gocore/v2/lib/tracing/client/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

endPointUrl := os.Getenv("ZIPKIN_BASE_URL")
// 务必遵循此规则，在链路追踪日志系统里面的应用名称规则
// 项目名称#环境 格式
appName := conf.ProjectName+"#"+utils.GetRunTime()
// 初始化tracer，只能初始化一次
// 第三个参数是采样率，0到1之间
_, err := tracing.InitZipkinTracer(appName, endPointUrl, 1)
if err != nil {
  panic(err)
}
// Set up a connection to the server peer.
conn, err := grpc.Dial(
  address,
  //... other options
  grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
  grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

// All future RPC activity involving `conn` will be automatically traced.
```

- grpc服务端使用方法

```go
import (
	tracing "github.com/sunmi-OS/gocore/v2/lib/tracing/client/zipkin-otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

endPointUrl := os.Getenv("ZIPKIN_BASE_URL")
// 务必遵循此规则，在链路追踪日志系统里面的应用名称规则
// 项目名称#环境 格式
appName := conf.ProjectName+"#"+utils.GetRunTime()
// 初始化tracer，一次运行中只能初始化一次
// 第三个参数是采样率，0到1之间
_, err := tracing.NewZipKinTracer(appName, endPointUrl, 1)
if err != nil {
  panic(err)
}
// ...

// Initialize the gRPC server.
s := grpc.NewServer(
    ... // other options
    grpc.UnaryInterceptor(
        otelgrpc.UnaryServerInterceptor()),
    grpc.StreamInterceptor(
        otelgrpc.StreamServerInterceptor()))

// All future RPC activity involving `s` will be automatically traced.
```

## gorm 

支持 gorm v2 链路追踪，使用方法如下：
修改 dal 层相关数据库的 mysql_client.go 文件，增加支持链路追踪的方法生成 db 对象

```go
package dal

import (
	"context"
	"github.com/sunmi-OS/gocore/v2/db/orm"
	tracing "github.com/sunmi-OS/gocore/v2/lib/tracing/gorm/otel"
	"gorm.io/gorm"
)

func OrmWithTracing(ctx context.Context) *gorm.DB {
	db := orm.GetORM("xxx")
	span, db := tracing.StartSpanWithCtx(ctx, db, 2)
	defer span.End()
	return db
}
```

数据库访问方法示例，需要多传入一个 context.Context 类型的参数

```go
func (m *User) GetById(ctx context.Context, id int64) (User, error) {
	var info User
	err := OrmWithTracing(ctx).Table(m.TableName()).Where("id = ?", id).Find(&info).Error
	return info, err
}
```

## Redis

在初始化 redis 的地方修改，修改 app/cmd/init.go 示例如下：

```go
import (
  "github.com/redis/go-redis/extra/redisotel/v9"
  "github.com/sunmi-OS/gocore/v2/db/redis"
  //...
)

// 初始化redis
func initRedis() {
	redis.NewRedis("xxx")
  
	rdb := redis.GetRedis("xxx")
	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}
	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}
}
```