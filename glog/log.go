package glog

import "github.com/sunmi-OS/gocore/v2/glog/zap"

func Info(args ...interface{}) {
	zap.Sugar.Info(args...)
}

func InfoF(format string, args ...interface{}) {
	zap.Sugar.Infof(format, args...)
}

func Debug(args ...interface{}) {
	zap.Sugar.Debug(args...)
}

func DebugF(format string, args ...interface{}) {
	zap.Sugar.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	zap.Sugar.Warn(args...)
}

func WarnF(format string, args ...interface{}) {
	zap.Sugar.Warnf(format, args...)
}

func Error(args ...interface{}) {
	zap.Sugar.Error(args...)
}

func ErrorF(format string, args ...interface{}) {
	zap.Sugar.Errorf(format, args...)
}
