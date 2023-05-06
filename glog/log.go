package glog

import (
	"sync"

	"github.com/sunmi-OS/gocore/v2/glog/logx"
	"github.com/sunmi-OS/gocore/v2/glog/zap"
)

var (
	Logger sync.Map
)

//  默认加入zap组件
func init() {
	Logger.Store("zap", &zap.Zap{})
}

// SetLogger设置日志打印实例,选择输出到文件,终端,阿里云日志等
func SetLogger(name string, logger logx.GLog) {
	Logger.Store(name, logger)
}

// DelLogger删除日志插件
func DelLogger(name string) {
	Logger.Delete(name)
}

func Info(args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Info(args...)
		return true
	})
}

func InfoF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).InfoF(format, args...)
		return true
	})
}

func Debug(args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Debug(args...)
		return true
	})
}

func DebugF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).DebugF(format, args...)
		return true
	})
}

func Warn(args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Warn(args...)
		return true
	})
}

func WarnF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).WarnF(format, args...)
		return true
	})
}

func Error(args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Error(args...)
		return true
	})
}

func ErrorF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).ErrorF(format, args...)
		return true
	})
}

func Fatal(args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Error(args...)
		return true
	})
}

func FatalF(format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).ErrorF(format, args...)
		return true
	})
}
