package utils

import (
	"fmt"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	var (
		count     int32 = 0
		execCount int32 = 0
	)
	// 重试5次，前三次失败，之后成功
	err := Retry(func() error {
		atomic.AddInt32(&execCount, 1)
		log.Println("exec:", execCount)
		if count < 3 {
			atomic.AddInt32(&count, 1)
			log.Println("retry again:", count)
			return fmt.Errorf("%d count retry finished", count)
		}
		log.Println("success:", execCount)
		return nil
	}, 5, 2*time.Second)
	if err != nil {
		log.Println(err)
	}
}
