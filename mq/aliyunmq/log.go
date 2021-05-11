package aliyunmq

import "github.com/apache/rocketmq-client-go/v2/rlog"

func LogInfo() {
	rlog.SetLogLevel("info")
}

func LogDebug() {
	rlog.SetLogLevel("debug")
}

func LogWarn() {
	rlog.SetLogLevel("warn")
}

func LogError() {
	rlog.SetLogLevel("error")
}
