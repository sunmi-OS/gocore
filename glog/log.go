package glog

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"github.com/sunmi-OS/gocore/v2/glog/logx"
	"github.com/sunmi-OS/gocore/v2/glog/zap"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	Logger sync.Map
	// statistic log err count
	metricErrCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "gocore",
			Name:      "log_error_total",
			Help:      "statistic log err count",
		}, []string{"source"})
)

const (
	_callerSkip      = "caller_skip"
	_defaultCallSkip = 2
)

// 默认加入zap组件
func init() {
	prometheus.MustRegister(metricErrCount)
	Logger.Store("zap", &zap.Zap{})
}

// SetLogger设置日志打印实例,选择输出到文件,终端,阿里云日志等
func SetLogger(name string, logger logx.GLog) {
	Logger.Store(name, logger)
}

// WithCallerSkipCtx 提供设置日志自定义source位置的入口
// 默认是3，当需要增加该值，可以设置获得更外层的
func WithCallerSkipCtx(ctx context.Context, skip int) context.Context {
	return context.WithValue(ctx, _callerSkip, skip)
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
	fn := funcName(_defaultCallSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Error(args...)
		return true
	})
}

func ErrorF(format string, args ...interface{}) {
	fn := funcName(_defaultCallSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).ErrorF(format, args...)
		return true
	})
}

func Fatal(args ...interface{}) {
	fn := funcName(_defaultCallSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).Fatal(args...)
		return true
	})
}

func FatalF(format string, args ...interface{}) {
	fn := funcName(_defaultCallSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).FatalF(format, args...)
		return true
	})
}

func InfoW(keyvals ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelInfo, context.TODO(), keyvals...)
		return true
	})
}

func DebugW(keyvals ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelDebug, context.TODO(), keyvals...)
		return true
	})
}

func WarnW(keyvals ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelWarn, context.TODO(), keyvals...)
		return true
	})
}

func ErrorW(keyvals ...interface{}) {
	fn := funcName(_defaultCallSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelError, context.TODO(), keyvals...)
		return true
	})
}

func FatalW(keyvals ...interface{}) {
	fn := funcName(_defaultCallSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelFatal, context.TODO(), keyvals...)
		return true
	})
}

func InfoC(ctx context.Context, format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelInfo, ctx, fmt.Sprintf(format, args...))
		return true
	})
}

func DebugC(ctx context.Context, format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelDebug, ctx, fmt.Sprintf(format, args...))
		return true
	})
}

func WarnC(ctx context.Context, format string, args ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelWarn, ctx, fmt.Sprintf(format, args...))
		return true
	})
}

func ErrorC(ctx context.Context, format string, args ...interface{}) {
	funcSkip := 0
	if value := ctx.Value(_callerSkip); value != nil {
		if i, ok := value.(int); ok {
			funcSkip = i
		}
	}
	fn := funcName(_defaultCallSkip + funcSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelError, ctx, fmt.Sprintf(format, args...))
		return true
	})
}

func FatalC(ctx context.Context, format string, args ...interface{}) {
	funcSkip := 0
	if value := ctx.Value(_callerSkip); value != nil {
		if i, ok := value.(int); ok {
			funcSkip = i
		}
	}
	fn := funcName(_defaultCallSkip + funcSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelFatal, ctx, fmt.Sprintf(format, args...))
		return true
	})
}

func InfoV(ctx context.Context, keyvals ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelInfo, ctx, keyvals...)
		return true
	})
}

func DebugV(ctx context.Context, keyvals ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelDebug, ctx, keyvals...)
		return true
	})
}

func WarnV(ctx context.Context, keyvals ...interface{}) {
	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelWarn, ctx, keyvals...)
		return true
	})
}

func ErrorV(ctx context.Context, keyvals ...interface{}) {
	funcSkip := 0
	if value := ctx.Value(_callerSkip); value != nil {
		if i, ok := value.(int); ok {
			funcSkip = i
		}
	}
	fn := funcName(_defaultCallSkip + funcSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelError, ctx, keyvals...)
		return true
	})
}

func FatalV(ctx context.Context, keyvals ...interface{}) {
	funcSkip := 0
	if value := ctx.Value(_callerSkip); value != nil {
		if i, ok := value.(int); ok {
			funcSkip = i
		}
	}
	fn := funcName(_defaultCallSkip + funcSkip)
	metricErrCount.WithLabelValues(fn).Inc()

	Logger.Range(func(k, v interface{}) bool {
		v.(logx.GLog).CommonLog(logx.LevelFatal, ctx, keyvals...)
		return true
	})
}

// funcName get func name.
func funcName(skip int) (name string) {
	if _, file, lineNo, ok := runtime.Caller(skip); ok {
		return file + ":" + strconv.Itoa(lineNo)
	}
	return "unknown:0"
}
