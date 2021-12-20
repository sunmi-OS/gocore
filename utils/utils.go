package utils

import (
	"fmt"
	"os"
	"time"
)

const (
	ReleaseEnv = "onl"
)

var releaseFlag = false //为true时表示线上环境

// GetDate 返回当前时间
func GetDate() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05")
}

// GetRunTime 获取当前系统环境
func GetRunTime() string {
	RunTime := os.Getenv("RUN_TIME")
	if RunTime == "" {
		fmt.Println("No RUN_TIME Can't start")
	}
	return RunTime
}

// 开启线上环境
func SetReleaseOn() error {
	releaseFlag = true
	return nil
}

// 如果是线上环境返回true
func IsRelease() bool {
	return releaseFlag || GetRunTime() == ReleaseEnv
}

// Either 返回一个存在的字符串
func Either(list ...string) string {
	for _, v := range list {
		if v != "" {
			return v
		}
	}
	return ""
}
