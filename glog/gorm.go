package glog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type dbLogger struct {
	SlowThreshold                       time.Duration // 慢 SQL 阈值
	IgnoreNotFoundError                 bool
	ParameterizedQueries                bool
	LogLevel                            logger.LogLevel
	traceStr, traceErrStr, traceWarnStr string
}

// NewDBLogger initialize db logger
func NewDBLogger(debug bool) logger.Interface {
	l := &dbLogger{
		SlowThreshold:        200 * time.Millisecond,
		IgnoreNotFoundError:  true,
		ParameterizedQueries: false,
		LogLevel:             logger.Warn,
		traceStr:             "[%.3fms] [rows:%v] %s",
		traceErrStr:          "[err=%+v] [%.3fms] [rows:%v] %s",
		traceWarnStr:         "[SLOW SQL >= %v] [%.3fms] [rows:%v] %s",
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
			"file_line", utils.FileWithLineNum(),
			"content", fmt.Sprintf(msg, data...),
		)
	}
}

// Warn print warn messages
func (l *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		WarnV(ctx,
			"kind", "SQL",
			"file_line", utils.FileWithLineNum(),
			"content", fmt.Sprintf(msg, data...),
		)
	}
}

// Error print error messages
func (l *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		ErrorV(ctx,
			"kind", "SQL",
			"file_line", utils.FileWithLineNum(),
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
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			ErrorV(ctx,
				"kind", "SQL",
				"file_line", utils.FileWithLineNum(),
				"content", fmt.Sprintf(l.traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql),
			)
		} else {
			ErrorV(ctx,
				"kind", "SQL",
				"file_line", utils.FileWithLineNum(),
				"content", fmt.Sprintf(l.traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql),
			)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		if rows == -1 {
			WarnV(ctx,
				"kind", "SQL",
				"file_line", utils.FileWithLineNum(),
				"content", fmt.Sprintf(l.traceWarnStr, l.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, "-", sql),
			)
		} else {
			WarnV(ctx,
				"kind", "SQL",
				"file_line", utils.FileWithLineNum(),
				"content", fmt.Sprintf(l.traceWarnStr, l.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, rows, sql),
			)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			InfoV(ctx,
				"kind", "SQL",
				"file_line", utils.FileWithLineNum(),
				"content", fmt.Sprintf(l.traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql),
			)
		} else {
			InfoV(ctx,
				"kind", "SQL",
				"file_line", utils.FileWithLineNum(),
				"content", fmt.Sprintf(l.traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql),
			)
		}
	}
}

// ParamsFilter filter params
func (l *dbLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
