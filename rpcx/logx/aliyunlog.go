package logx

import (
	"fmt"

	"github.com/sunmi-OS/gocore/v2/glog/aliyunlog"
)

type AliyunLogger struct {
}

func (s *AliyunLogger) Info(key string, content map[string]string) {
	err := aliyunlog.Info(key, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}

func (s *AliyunLogger) Error(key string, content map[string]string) {
	err := aliyunlog.Error(key, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}
