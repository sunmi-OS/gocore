package retry

import (
	"context"
	"time"
)

// Retry 重试
func Retry(ctx context.Context, fn func(ctx context.Context) error, attempts int, sleep time.Duration) (err error) {
	threshold := attempts - 1
	for i := 0; i < attempts; i++ {
		err = fn(ctx)
		if err == nil || i >= threshold {
			return
		}
		time.Sleep(sleep)
	}
	return
}
