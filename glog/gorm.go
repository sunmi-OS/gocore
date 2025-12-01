package glog

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type dbLogger struct {
	SlowThreshold                           time.Duration // 慢 SQL 阈值
	IgnoreNotFoundError                     bool
	ParameterizedQueries                    bool
	LogLevel                                logger.LogLevel
	traceInfoStr, traceErrStr, traceWarnStr string
}

// NewDBLogger initialize db logger
func NewDBLogger(debug bool, slowThreshold time.Duration) logger.Interface {
	l := &dbLogger{
		SlowThreshold:        slowThreshold,
		IgnoreNotFoundError:  true,
		ParameterizedQueries: false,
		LogLevel:             logger.Warn,
		traceInfoStr:         "[%.3fms] [rows:%s] %s",
		traceErrStr:          "[err=%+v] [%.3fms] [rows:%s] %s",
		traceWarnStr:         "[slow_sql >= %v] [%.3fms] [rows:%s] %s",
	}
	if debug {
		l.LogLevel = logger.Info
	}
	return l
}

// LogMode log mode
func (l *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlog := *l
	newlog.LogLevel = level
	return &newlog
}

// Info print info
func (l *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		InfoV(ctx,
			"kind", "SQL",
			"file_line", getCallerFileWithLineNum(),
			"content", fmt.Sprintf(msg, data...),
		)
	}
}

// Warn print warn messages
func (l *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		WarnV(ctx,
			"kind", "SQL",
			"file_line", getCallerFileWithLineNum(),
			"content", fmt.Sprintf(msg, data...),
		)
	}
}

// Error print error messages
func (l *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		ErrorV(ctx,
			"kind", "SQL",
			"file_line", getCallerFileWithLineNum(),
			"content", fmt.Sprintf(msg, data...),
		)
	}
}

// Trace print sql message
func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	// trace log
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreNotFoundError):
		l.logTrace(ctx, logger.Error, elapsed, rows, sql, err)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.logTrace(ctx, logger.Warn, elapsed, rows, sql, nil)
	case l.LogLevel == logger.Info:
		l.logTrace(ctx, logger.Info, elapsed, rows, sql, nil)
	}
}

func (l *dbLogger) logTrace(ctx context.Context, level logger.LogLevel, elapsed time.Duration, rows int64, sql string, err error) {
	rowStr := "-"
	if rows >= 0 {
		rowStr = strconv.FormatInt(rows, 10)
	}
	var (
		logFn   func(ctx context.Context, keyvals ...interface{})
		content string
	)
	switch level {
	case logger.Info:
		logFn = InfoV
		content = fmt.Sprintf(l.traceInfoStr, float64(elapsed.Nanoseconds())/1e6, rowStr, sql)
	case logger.Warn:
		logFn = WarnV
		content = fmt.Sprintf(l.traceWarnStr, l.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, rowStr, sql)
	case logger.Error:
		logFn = ErrorV
		content = fmt.Sprintf(l.traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rowStr, sql)
	default:
		return
	}
	logFn(ctx,
		"kind", "SQL",
		"file_line", getCallerFileWithLineNum(),
		"content", content,
	)
}

// ParamsFilter filter params
func (l *dbLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

// getCallerFileWithLineNum 获取调用者的文件名和行号，跳过 gorm 和 glog 内部调用
func getCallerFileWithLineNum() string {
	// 从调用栈中查找第一个不是 gorm 和 glog 包的调用位置
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			// 跳过 gorm 和 glog 包
			if !strings.Contains(file, "gorm.io/gorm") &&
				!strings.Contains(file, "/glog/") &&
				!strings.HasSuffix(file, "/glog/gorm.go") {
				return fmt.Sprintf("%s:%d", file, line)
			}
		}
	}
	// 如果找不到，使用 gorm 的默认方法
	return utils.FileWithLineNum()
}
