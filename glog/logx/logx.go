package logx

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
}
