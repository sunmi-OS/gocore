// Package utils provides retry times functions.
// Author: Jerry
package utils

import (
	"time"
)

// Retry 重试 func 最大次数，间隔
func Retry(fc func() error, maxRetries int, interval time.Duration) (err error) {
	for i := 1; i <= maxRetries; i++ {
		if err = fc(); err != nil {
			time.Sleep(interval)
			continue
		}
		return nil
	}
	return err
}
