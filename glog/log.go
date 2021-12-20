package glog

import (
	"github.com/sunmi-OS/gocore/v2/glog/logx"
	"github.com/sunmi-OS/gocore/v2/glog/zap"
)

var (
	Logger logx.GLog
)

func init() {
	// 默认使用zap打印日志
	SetLogger(&zap.Zap{})
}

// 设置日志打印实例,选择输出到文件,终端,阿里云日志等
func SetLogger(logger logx.GLog) {
	Logger = logger
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func InfoF(format string, args ...interface{}) {
	Logger.InfoF(format, args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func DebugF(format string, args ...interface{}) {
	Logger.DebugF(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func WarnF(format string, args ...interface{}) {
	Logger.WarnF(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func ErrorF(format string, args ...interface{}) {
	Logger.ErrorF(format, args...)
}
