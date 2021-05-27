package logx

import (
	"fmt"
	aliyunlog2 "github.com/sunmi-OS/gocore/v2/utils/aliyunlog"
)

type AliyunLogger struct {
}

func (s *AliyunLogger) Info(key string, content map[string]string) {
	err := aliyunlog2.Info(key, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}

func (s *AliyunLogger) Error(key string, content map[string]string) {
	err := aliyunlog2.Error(key, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}
