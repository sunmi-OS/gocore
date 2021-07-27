package utils

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog/xlog"
)

func TestRetry(t *testing.T) {
	var count int32 = 0

	err := Retry(func() error {
		if count < 3 {
			atomic.AddInt32(&count, 1)
			return fmt.Errorf("%d count retry finished", count)
		}
		return nil
	}, 5, 2*time.Second)
	if err != nil {
		xlog.Error(err)
	}
}
