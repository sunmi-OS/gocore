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
	if logger.Sugar != nil {
		logger.Sugar.Info(args...)
		return
	}
	infoLog.logOut(nil, args...)
}

func Infof(format string, args ...interface{}) {
	if logger.Sugar != nil {
		logger.Sugar.Infof(format, args...)
		return
	}
	infoLog.logOut(&format, args...)
}

func Debug(args ...interface{}) {
	if logger.Sugar != nil {
		logger.Sugar.Debug(args...)
		return
	}
	debugLog.logOut(nil, args...)
}

func Debugf(format string, args ...interface{}) {
	if logger.Sugar != nil {
		logger.Sugar.Debugf(format, args...)
		return
	}
	debugLog.logOut(&format, args...)
}

func Warn(args ...interface{}) {
	if logger.Sugar != nil {
		logger.Sugar.Warn(args...)
		return
	}
	warnLog.logOut(nil, args...)
}

func Warnf(format string, args ...interface{}) {
	if logger.Sugar != nil {
		logger.Sugar.Warnf(format, args...)
		return
	}
	warnLog.logOut(&format, args...)
}

func Error(args ...interface{}) {
	if logger.Sugar != nil {
		logger.Sugar.Error(args...)
		return
	}
	errLog.logOut(nil, args...)
}

func Errorf(format string, args ...interface{}) {
	if logger.Sugar != nil {
		logger.Sugar.Errorf(format, args...)
		return
	}
	errLog.logOut(&format, args...)
}
