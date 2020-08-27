package log

import (
	"testing"
	"time"

	"github.com/sunmi-OS/gocore/xlog"
)

func TestInitLogger(t *testing.T) {
	//InitLogger()
	//Logger.Info("sunmi")
	//Sugar.Info("sunmi")

	now := time.Now()
	// 计算下一个零点
	next := now.Add(time.Hour * 24)
	xlog.Info(next)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	time.NewTimer(next.Sub(now))
	xlog.Info(next)

}
