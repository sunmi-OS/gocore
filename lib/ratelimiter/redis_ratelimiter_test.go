package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// 创建测试用的 Redis 客户端
func setupTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// 确保连接成功
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 清理测试数据
	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("Failed to flush Redis DB: %v", err)
	}

	return client
}

func TestNewRedisRateLimiter(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	tests := []struct {
		name    string
		config  RedisConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: RedisConfig{
				Rate:   "10-S",
				Prefix: "test",
			},
			wantErr: false,
		},
		{
			name: "empty prefix",
			config: RedisConfig{
				Rate: "10-S",
			},
			wantErr: false,
		},
		{
			name: "invalid rate format",
			config: RedisConfig{
				Rate:   "invalid",
				Prefix: "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter, err := NewRedisRateLimiter(client, tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if limiter == nil {
				t.Error("Expected limiter to be non-nil")
			}
		})
	}
}

func TestRedisRateLimiter_Get(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	// 创建限流器：每秒2次请求
	config := RedisConfig{
		Rate:   "2-S",
		Prefix: "test",
	}
	limiter, err := NewRedisRateLimiter(client, config)
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}
	if limiter == nil {
		t.Fatal("Expected limiter to be non-nil")
	}

	ctx := context.Background()
	key := "test-key"

	// 测试第一次请求
	context, err := limiter.Get(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error on first request: %v", err)
		return
	}
	if context.Reached {
		t.Error("Expected first request to not be reached")
	}
	if context.Limit != 2 {
		t.Errorf("Expected limit to be 2, got %d", context.Limit)
	}
	if context.Remaining != 1 {
		t.Errorf("Expected remaining to be 1, got %d", context.Remaining)
	}

	// 测试第二次请求
	context, err = limiter.Get(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error on second request: %v", err)
		return
	}
	if context.Reached {
		t.Error("Expected second request to not be reached")
	}
	if context.Limit != 2 {
		t.Errorf("Expected limit to be 2, got %d", context.Limit)
	}
	if context.Remaining != 0 {
		t.Errorf("Expected remaining to be 0, got %d", context.Remaining)
	}

	// 测试第三次请求（应该被限流）
	context, err = limiter.Get(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error on third request: %v", err)
		return
	}
	if !context.Reached {
		t.Error("Expected third request to be reached")
	}
	if context.Limit != 2 {
		t.Errorf("Expected limit to be 2, got %d", context.Limit)
	}
	if context.Remaining != 0 {
		t.Errorf("Expected remaining to be 0, got %d", context.Remaining)
	}

	// 等待一秒后重置
	time.Sleep(time.Second)

	// 测试重置后的请求
	context, err = limiter.Get(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error after reset: %v", err)
		return
	}
	if context.Reached {
		t.Error("Expected request after reset to not be reached")
	}
	if context.Limit != 2 {
		t.Errorf("Expected limit to be 2, got %d", context.Limit)
	}
	if context.Remaining != 1 {
		t.Errorf("Expected remaining to be 1, got %d", context.Remaining)
	}
}

func TestRedisRateLimiter_Concurrent(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	// 创建限流器：每秒10次请求
	config := RedisConfig{
		Rate:   "10-S",
		Prefix: "test",
	}
	limiter, err := NewRedisRateLimiter(client, config)
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}
	if limiter == nil {
		t.Fatal("Expected limiter to be non-nil")
	}

	ctx := context.Background()
	key := "test-key"

	// 并发测试
	concurrent := 15
	results := make(chan bool, concurrent)

	for i := 0; i < concurrent; i++ {
		go func() {
			context, err := limiter.Get(ctx, key)
			if err != nil {
				results <- false
				return
			}
			results <- !context.Reached
		}()
	}

	// 收集结果
	successCount := 0
	for i := 0; i < concurrent; i++ {
		if <-results {
			successCount++
		}
	}

	// 验证结果：应该只有10个请求成功
	if successCount != 10 {
		t.Errorf("Expected 10 successful requests, got %d", successCount)
	}
}

func TestRedisRateLimiter_DifferentKeys(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	// 创建限流器：每秒2次请求
	config := RedisConfig{
		Rate:   "2-S",
		Prefix: "test",
	}
	limiter, err := NewRedisRateLimiter(client, config)
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}
	if limiter == nil {
		t.Fatal("Expected limiter to be non-nil")
	}

	ctx := context.Background()
	key1 := "key1"
	key2 := "key2"

	// 测试第一个key
	context, err := limiter.Get(ctx, key1)
	if err != nil {
		t.Errorf("Unexpected error on first key: %v", err)
		return
	}
	if context.Reached {
		t.Error("Expected first key request to not be reached")
	}
	if context.Remaining != 1 {
		t.Errorf("Expected remaining to be 1, got %d", context.Remaining)
	}

	// 测试第二个key（应该不受第一个key的影响）
	context, err = limiter.Get(ctx, key2)
	if err != nil {
		t.Errorf("Unexpected error on second key: %v", err)
		return
	}
	if context.Reached {
		t.Error("Expected second key request to not be reached")
	}
	if context.Remaining != 1 {
		t.Errorf("Expected remaining to be 1, got %d", context.Remaining)
	}
}
