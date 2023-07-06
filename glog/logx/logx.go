package logx

import "context"

type Level int8

const (
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type LogType int64

const (
	LogTypeZap LogType = iota + 1
	LogTypeSls
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return ""
	}
}

type GLog interface {
	Info(args ...interface{})
	InfoF(format string, args ...interface{})

	Debug(args ...interface{})
	DebugF(format string, args ...interface{})

	Warn(args ...interface{})
	WarnF(format string, args ...interface{})

	Error(args ...interface{})
	ErrorF(format string, args ...interface{})

	Fatal(args ...interface{})
	FatalF(format string, args ...interface{})

	CommonLog(level Level, ctx context.Context, keyvals ...interface{}) error
}
