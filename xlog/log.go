package xlog

var (
	debugLog xLogger = &DebugLogger{}
	infoLog  xLogger = &InfoLogger{}
	warnLog  xLogger = &WarnLogger{}
	errLog   xLogger = &ErrorLogger{}
)

type xLogger interface {
	logOut(format *string, args ...interface{})
}

func Info(args ...interface{}) {
	infoLog.logOut(nil, args...)
}

func Infof(format string, args ...interface{}) {
	infoLog.logOut(&format, args...)
}

func Debug(args ...interface{}) {
	debugLog.logOut(nil, args...)
}

func Debugf(format string, args ...interface{}) {
	debugLog.logOut(&format, args...)
}

func Warn(args ...interface{}) {
	warnLog.logOut(nil, args...)
}

func Warnf(format string, args ...interface{}) {
	warnLog.logOut(&format, args...)
}

func Error(args ...interface{}) {
	errLog.logOut(nil, args...)
}

func Errorf(format string, args ...interface{}) {
	errLog.logOut(&format, args...)
}
