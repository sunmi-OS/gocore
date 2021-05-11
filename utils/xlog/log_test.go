package xlog

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	// default log
	Info("info")
	Debug("debug")
	Warn("warning")
	Error("error")

	fmt.Println()

	// zap log
	Zap().Infof("%+v", struct {
		Name string
		Age  int
	}{
		Name: "Jerry",
		Age:  18,
	})
	Zap().Debug("zap debug")
	Zap().Warn("zap warn")
	Zap().Error("zap error")
}
