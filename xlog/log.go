package xlog

func Info(args ...interface{}) {
	logger.Sugar.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Sugar.Infof(format, args...)
}

func Debug(args ...interface{}) {
	logger.Sugar.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Sugar.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Sugar.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Sugar.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logger.Sugar.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Sugar.Errorf(format, args...)
}
