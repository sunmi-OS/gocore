package utils

import (
	"os"
	"time"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// 返回当前时间
func GetDate() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05")
}

// 获取当前系统环境
func GetRunTime() string {
	//获取系统环境变量
	RUN_TIME := os.Getenv("RUN_TIME")
	if RUN_TIME == "" {
		fmt.Println("No RUN_TIME Can't start")
	}
	return RUN_TIME
}

// MD5 加密字符串
func GetMD5(plainText string) string {
	h := md5.New()
	h.Write([]byte(plainText))
	return hex.EncodeToString(h.Sum(nil))
}
