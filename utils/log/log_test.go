package log

import (
	xlog2 "github.com/sunmi-OS/gocore/v2/utils/xlog"
	"testing"
	"time"
)

func TestInitLogger(t *testing.T) {
	//InitLogger()
	//Logger.Info("sunmi")
	//Sugar.Info("sunmi")

	now := time.Now()
	// 计算下一个零点
	next := now.Add(time.Hour * 24)
	xlog2.Info(next)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	time.NewTimer(next.Sub(now))
	xlog2.Info(next)

}
