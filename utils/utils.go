package utils

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	LocalEnv   = "local"
	DevEnv     = "dev"
	TestEnv    = "test"
	UatEnv     = "uat"
	PreEnv     = "pre"
	ReleaseEnv = "onl"
)

var (
	releaseFlag = false // 为true时表示线上环境
	runTime     string
	appName     string
	zone        string
	hostname    string
)

const TimeFormat = "2006-01-02T15:04:05.000Z0700"

// GetDate 返回当前时间
func GetDate() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05")
}

// GetRunTime 获取当前系统环境
func GetRunTime() string {
	if runTime != "" {
		return runTime
	}
	runTime = os.Getenv("RUN_TIME")
	if runTime == "" {
		fmt.Println("No RUN_TIME Can't start")
	}
	return runTime
}

func GetAppName() string {
	if appName != "" {
		return appName
	}
	appName = os.Getenv("APP_NAME")
	if appName == "" {
		fmt.Println("No APP_NAME Set")
		appName = "UnknownAppName"
	}
	return appName
}

func GetZone() string {
	if zone != "" {
		return zone
	}
	zone = os.Getenv("ZONE")
	if zone == "" {
		fmt.Println("No ZONE Set")
	}
	return zone
}

func GetHostname() string {
	if hostname != "" {
		return hostname
	}
	hostname = os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	return hostname
}

// OnRelease 开启线上环境
func OnRelease() {
	releaseFlag = true
}

// IsRelease 如果是线上环境返回true
func IsRelease() bool {
	return releaseFlag || GetRunTime() == ReleaseEnv
}

func IsLocal() bool {
	return GetRunTime() == LocalEnv
}

func IsDev() bool {
	return GetRunTime() == DevEnv
}

func IsTest() bool {
	return GetRunTime() == TestEnv
}

func IsUat() bool {
	return GetRunTime() == UatEnv
}

func IsPre() bool {
	return GetRunTime() == PreEnv
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

// GetAccesslogPath accesslog路径
func GetAccesslogPath() string {
	var path string
	switch runtime.GOOS {
	case "windows":
		path = "./logs/access.log"
	case "darwin":
		path = "./logs/access.log"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		path = "/data/logs/access.log"
	}
	return path
}
