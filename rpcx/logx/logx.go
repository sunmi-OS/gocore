package logx

import (
	"encoding/json"
	"fmt"
	"time"
)

var LoggerObj Logger = &DefaultLogger{}

func SetLogger(l Logger) {
	LoggerObj = l
}

type Logger interface {
	Info(key string, content map[string]string)
	Error(key string, content map[string]string)
}

type DefaultLogger struct {
}

type LogData struct {
	Level   string            `json:"level"`
	Time    string            `json:"time"`
	Key     string            `json:"key"`
	Content map[string]string `json:"content"`
}

func (s *DefaultLogger) Info(key string, content map[string]string) {
	data := LogData{
		Level:   "info",
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Key:     key,
		Content: content,
	}
	dataByte, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	fmt.Println(string(dataByte))
}

func (s *DefaultLogger) Error(key string, content map[string]string) {
	data := LogData{
		Level:   "error",
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Key:     key,
		Content: content,
	}
	dataByte, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	fmt.Println(string(dataByte))
}
