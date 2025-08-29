package ratelimiter

import (
	"context"
	"fmt"
	"regexp"

	goRedis "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

// Validate 验证配置是否有效
func (c RedisConfig) Validate() error {
	if c.Rate == "" {
		return fmt.Errorf("rate cannot be empty")
	}

	// 验证速率格式
	rateRegex := regexp.MustCompile(`^\d+-[SMHD]$`)
	if !rateRegex.MatchString(c.Rate) {
		return fmt.Errorf("invalid rate format: %s, expected format: <number>-<period> (e.g. 10-S)", c.Rate)
	}

	// 验证前缀
	if c.Prefix == "" {
		c.Prefix = "redisLimiter"
	}
	if len(c.Prefix) > 50 {
		return fmt.Errorf("prefix too long, maximum length is 50 characters")
	}

	return nil
}

// RedisConfig 限流器配置
type RedisConfig struct {
	// 限流速率，例如 "10-S" 表示每秒10次
	// You can also use the simplified format "<limit>-<period>"", with the given
	// periods:
	//
	// * "S": second
	// * "M": minute
	// * "H": hour
	// * "D": day
	//
	// Examples:
	//
	// * 5 reqs/second: "5-S"
	// * 10 reqs/minute: "10-M"
	// * 1000 reqs/hour: "1000-H"
	// * 2000 reqs/day: "2000-D"
	Rate string
	// 前缀
	Prefix string
}

// RedisRateLimiter 限流器接口
type RedisRateLimiter struct {
	limiter *limiter.Limiter
	config  RedisConfig
}

// NewRedisRateLimiter 创建新的限流器
func NewRedisRateLimiter(redisClient *goRedis.Client, config RedisConfig) (*RedisRateLimiter, error) {
	// 添加参数验证
	if redisClient == nil {
		return nil, fmt.Errorf("redis client cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 创建Redis存储
	store, err := redis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix: config.Prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create redis store: %w", err)
	}

	// 创建限流器
	rate, err := limiter.NewRateFromFormatted(config.Rate)
	if err != nil {
		return nil, fmt.Errorf("invalid rate format: %w", err)
	}

	limiter := limiter.New(store, rate)

	return &RedisRateLimiter{
		limiter: limiter,
		config:  config,
	}, nil
}

// Get 检查是否被限流
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - key: 限流对象的唯一标识符
//
// 返回值:
//   - limiter.Context: 包含限流状态的结构体，具体字段如下：
//   - Reached: 是否达到限流阈值（true表示已限流，false表示未限流）
//   - Limit: 限流器的总配额（例如：配置"2-S"时为2,表示每秒2个请求）
//   - Remaining: 当前时间窗口内的剩余可用请求数(比如在测试用例中：第一次请求后，Remaining = 1;第二次请求后，Remaining = 0; 第三次请求时因为已经到达限制，Remaining 保持为 0)
//   - error: 操作过程中发生的错误，如果操作成功则为nil
func (r *RedisRateLimiter) Get(ctx context.Context, key string) (limiter.Context, error) {
	return r.limiter.Get(ctx, key)
}
