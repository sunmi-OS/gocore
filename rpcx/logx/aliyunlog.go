package logx

import (
	"fmt"

	"github.com/sunmi-OS/gocore/aliyunlog"
)

type aliyunLogger struct {
}

func (s *aliyunLogger) Info(key string, content map[string]string) {
	err := aliyunlog.Info(key, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}

func (s *aliyunLogger) Error(key string, content map[string]string) {
	err := aliyunlog.Info(key, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}
