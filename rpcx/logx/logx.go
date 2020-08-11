package logx

import (
	"fmt"
	"gocore/aliyunlog"
)

const timeFormat = "2006-01-02 15:04:05"

var LoggerObj Logger

func SetLogger(l Logger) {
	LoggerObj = l
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type defaultLogger struct {
}

func (s *defaultLogger) Info(key string, content map[string]string) {
	err := aliyunlog.Info(key, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}

func (s *defaultLogger) Error(key string, content map[string]string) {
	err := aliyunlog.Info(key, content)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
}
