package glog

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sunmi-OS/gocore/v2/glog/zap"
)

func TestLog(t *testing.T) {
	s := struct {
		Name string
		Age  int
	}{
		Name: "Jerry",
		Age:  18,
	}
	// zap log
	InfoF("%+v", s)
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")
	ErrorF("s.dao.PartnerById(%d),err:%+v", 10086, errors.New("不存在此id"))
	ErrorF("s.dao.CreateOrder(%+v),err:%+v", s, errors.New("创建订单失败"))
	ErrorW("key", "err value", "key3", "value3")
	ErrorW("key", "err value", "key3")
	//Fatal("zap fatal")

	fmt.Println("")

	zap.SetLogLevel(zap.LogLevelWarn)
	InfoF("%+v", s)
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")
	ErrorW("zap", "error")
	//Fatal("zap fatal")

	fmt.Println("")

	zap.InitFileLog()
	Debug("zap debug")
	Warn("zap warn")
	Error("zap error")
	ErrorW("zap", "error")
	//Fatal("zap fatal")

	//
	ctx := context.Background()
	InfoV(ctx, "key", "value")
	InfoV(ctx, "key", "value", "key2", "value2", "key3")
	InfoC(ctx, "format: %v", 12345)
	InfoW("key", "value", "key2", "value2")
	InfoW("key", "value", "key3")
}
