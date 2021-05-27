package retry

import (
	"fmt"
	xlog2 "github.com/sunmi-OS/gocore/v2/utils/xlog"
	"sync/atomic"
	"testing"
	"time"
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
		xlog2.Error(err)
	}
}
