# Rate Limiter

一个基于Redis的分布式限流器库，提供简单易用的API接口来控制请求频率。

## 特性

- 🚀 **高性能**: 基于Redis实现，支持分布式部署
- 🔧 **灵活配置**: 支持秒、分钟、小时、天级别的限流配置
- 🛡️ **参数验证**: 内置配置验证，确保参数正确性
- 📊 **详细统计**: 提供限流状态、剩余配额等详细信息
- 🎯 **简单易用**: 简洁的API设计，快速上手

## 安装

```bash
go get github.com/sunmi-OS/gocore/v2/lib/ratelimiter
```

## 快速开始

### 1. 创建Redis客户端

```go
import (
    "github.com/redis/go-redis/v9"
    "github.com/sunmi-OS/gocore/v2/lib/ratelimiter"
)

// 创建Redis客户端
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    DB:   0,
})
```

### 2. 配置限流器

```go
// 配置限流规则：每秒10次请求
config := ratelimiter.RedisConfig{
    Rate:   "10-S",  // 每秒10次
    Prefix: "api",   // Redis键前缀
}
```

### 3. 创建限流器实例

```go
limiter, err := ratelimiter.NewRedisRateLimiter(redisClient, config)
if err != nil {
    log.Fatal("Failed to create rate limiter:", err)
}
```

### 4. 使用限流器

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 使用用户ID作为限流键
    userID := getUserID(r)
    
    // 检查是否被限流
    ctx, err := limiter.Get(r.Context(), userID)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    
    if ctx.Reached {
        // 请求被限流
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }
    
    // 请求正常处理
    // ... 业务逻辑
}
```

## 配置说明

### 限流速率格式

限流速率使用 `<limit>-<period>` 格式：

- **S**: 秒 (Second)
- **M**: 分钟 (Minute)  
- **H**: 小时 (Hour)
- **D**: 天 (Day)

#### 示例

```go
// 每秒5次请求
config := ratelimiter.RedisConfig{
    Rate: "5-S",
}

// 每分钟100次请求
config := ratelimiter.RedisConfig{
    Rate: "100-M",
}

// 每小时1000次请求
config := ratelimiter.RedisConfig{
    Rate: "1000-H",
}

// 每天2000次请求
config := ratelimiter.RedisConfig{
    Rate: "2000-D",
}
```

### 配置参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| Rate | string | 是 | - | 限流速率，格式：`<limit>-<period>` |
| Prefix | string | 否 | "redisLimiter" | Redis键前缀，最大长度50字符 |

## API 参考

### RedisRateLimiter

#### NewRedisRateLimiter(redisClient, config)

创建新的限流器实例。

**参数:**
- `redisClient`: Redis客户端实例
- `config`: 限流器配置

**返回值:**
- `*RedisRateLimiter`: 限流器实例
- `error`: 错误信息

#### Get(ctx, key)

检查指定键是否被限流。

**参数:**
- `ctx`: 上下文对象
- `key`: 限流对象的唯一标识符

**返回值:**
- `limiter.Context`: 限流状态信息
- `error`: 错误信息

### limiter.Context 结构

```go
type Context struct {
    Reached   bool      // 是否达到限流阈值
    Limit     int64     // 限流器的总配额
    Remaining int64     // 当前时间窗口内的剩余可用请求数
    Reset     time.Time // 重置时间
    RetryAfter time.Duration // 重试等待时间
}
```

## 使用场景

- **API限流**: 防止API被恶意调用或过载
- **用户行为控制**: 限制用户操作频率
- **资源保护**: 保护数据库、缓存等资源
- **DDoS防护**: 防止分布式拒绝服务攻击

## 完整示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/redis/go-redis/v9"
    "github.com/sunmi-OS/gocore/v2/lib/ratelimiter"
)

func main() {
    // 创建Redis客户端
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   0,
    })
    
    // 创建限流器：每秒2次请求
    config := ratelimiter.RedisConfig{
        Rate:   "2-S",
        Prefix: "api_limiter",
    }
    
    limiter, err := ratelimiter.NewRedisRateLimiter(redisClient, config)
    if err != nil {
        log.Fatal("Failed to create rate limiter:", err)
    }
    
    // 模拟多个请求
    for i := 1; i <= 5; i++ {
        ctx, err := limiter.Get(context.Background(), "user123")
        if err != nil {
            log.Printf("Request %d: Error: %v", i, err)
            continue
        }
        
        if ctx.Reached {
            log.Printf("Request %d: Rate limited! Remaining: %d", i, ctx.Remaining)
        } else {
            log.Printf("Request %d: Allowed. Remaining: %d", i, ctx.Remaining)
        }
        
        time.Sleep(200 * time.Millisecond)
    }
}
```

## 测试

运行测试用例：

```bash
cd ratelimiter
go test -v
```

**注意**: 测试需要本地Redis服务运行在 `localhost:6379`

## 依赖

- `github.com/redis/go-redis/v9`: Redis客户端
- `github.com/ulule/limiter/v3`: 限流算法实现


## 贡献

欢迎提交Issue和Pull Request来改进这个库。
