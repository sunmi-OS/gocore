package utils

import (
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

type X map[string]any

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
	dcName      string
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
	return runTime
}

func GetDcName() string {
	if dcName != "" {
		return dcName
	}
	dcName = os.Getenv("DC_NAME")
	return dcName
}

func GetAppName() string {
	if appName != "" {
		return appName
	}
	appName = os.Getenv("APP_NAME")
	if appName == "" {
		appName = "UnknownAppName"
	}
	return appName
}

func GetZone() string {
	if zone != "" {
		return zone
	}
	zone = os.Getenv("ZONE")
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

// IsGlocal Check whether it is glocal
func IsGlocal() bool {
	dc := GetDcName()
	if dc == "CN" || dc == "" {
		return false
	}
	return true
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

// GetAccesslogPath accesslog path
func GetAccesslogPath() string {
	var path string
	path = os.Getenv("ACCESS_LOG_PATH")
	if path != "" {
		return path
	}
	appName = GetAppName()
	switch runtime.GOOS {
	case "windows":
		path = "./logs/" + appName + ".log"
	case "darwin":
		path = "./logs/" + appName + ".log"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		path = "/data/logs/" + appName + ".log"
	}
	return path
}

func LogContentUnmarshal(content string) interface{} {
	var value map[string]interface{}
	if strings.HasPrefix(content, "{") && sonic.UnmarshalString(content, &value) == nil {
		return value
	}
	return content
}

type Step struct {
	Head int
	Tail int
}

// Steps calculates the steps.
// example:
//
//	arr := []int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20}
//	for _, step := range golib.Steps(len(arr), 10) {
//		cur := arr[step.Head:step.Tail]
//		// todo: do something
//	}
func Steps(total, step int) (steps []Step) {
	steps = make([]Step, 0, int(math.Ceil(float64(total)/float64(step))))
	for i := 0; i < total; i++ {
		if i%step == 0 {
			head := i
			tail := head + step
			if tail > total {
				tail = total
			}
			steps = append(steps, Step{Head: head, Tail: tail})
		}
	}
	return steps
}

// IsUniqueDuplicateError 判断是否「唯一索引冲突」错误
func IsUniqueDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	for _, s := range []string{
		"Duplicate entry",            // MySQL
		"violates unique constraint", // Postgres
		"UNIQUE constraint failed",   // SQLite
	} {
		if strings.Contains(err.Error(), s) {
			return true
		}
	}
	return false
}
