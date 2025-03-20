package xlog

import (
	"runtime"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	debugLog xLogger = &DebugLogger{}
	infoLog  xLogger = &InfoLogger{}
	warnLog  xLogger = &WarnLogger{}
	errLog   xLogger = &ErrorLogger{}

	// statistic log err count
	metricErrCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gocore_1_5",
			Name:      "log_error_total",
			Help:      "statistic log err count",
		}, []string{"source"})
)

const (
	_callerSkip      = "caller_skip"
	_defaultCallSkip = 2
)

func init() {
	prometheus.MustRegister(metricErrCount)
}

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
	fn := funcName(_defaultCallSkip)
	metricErrCount.WithLabelValues(fn).Inc()
	errLog.logOut(nil, args...)
}

func Errorf(format string, args ...interface{}) {
	fn := funcName(_defaultCallSkip)
	metricErrCount.WithLabelValues(fn).Inc()
	errLog.logOut(&format, args...)
}

// funcName get func name.
func funcName(skip int) (name string) {
	if _, file, lineNo, ok := runtime.Caller(skip); ok {
		return file + ":" + strconv.Itoa(lineNo)
	}
	return "unknown:0"
}
