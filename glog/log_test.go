package glog

import (
	"testing"

	"github.com/sunmi-OS/gocore/v2/glog/zap"
)

func TestLog(t *testing.T) {
	// zap log
	InfoF("%+v", struct {
		Name string
		Age  int
	}{
		Name: "Jerry",
		Age:  18,
	})
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")

	zap.SetLocLevel("error")
	InfoF("%+v", struct {
		Name string
		Age  int
	}{
		Name: "Jerry",
		Age:  18,
	})
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")

	zap.InitFileLog()
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")

}
