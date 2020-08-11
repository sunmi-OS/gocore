package logx

import (
	"fmt"
	"time"
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
	fmt.Printf("%s:%s:%#v\n", time.Now().Format("2006-01-02 15:04:05"), key, content)
}

func (s *defaultLogger) Error(key string, content map[string]string) {
	fmt.Printf("%s:%s:%#v\n", time.Now().Format("2006-01-02 15:04:05"), key, content)
}
