package logx

import (
	"fmt"
	"time"
)

type aliyunLogger struct {
}

func (s *aliyunLogger) Info(key string, content map[string]string) {
	fmt.Printf("%s:%s:%#v\n", time.Now().Format("2006-01-02 15:04:05"), key, content)
}

func (s *aliyunLogger) Error(key string, content map[string]string) {
	fmt.Printf("%s:%s:%#v\n", time.Now().Format("2006-01-02 15:04:05"), key, content)
}
