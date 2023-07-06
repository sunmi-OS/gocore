package logx

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sunmi-OS/gocore/v2/utils"
	"go.opentelemetry.io/otel/trace"
)

type Valuer func(ctx context.Context) interface{}

var zapPrefixKVs = []interface{}{
	"appname", utils.GetAppName(),
	"runtime", utils.GetRunTime(),
	"hostname", utils.GetHostname(),
}

var slsPrefixKVs = []interface{}{
	"appname", utils.GetAppName(),
	"runtime", utils.GetRunTime(),
	"hostname", utils.GetHostname(),
	"ts", Timestamp(utils.TimeFormat),
	"caller", Caller(6),
}

// AppendKVs 可以进行自定义
func AppendKVs(keyvals ...interface{}) {
	zapPrefixKVs = append(zapPrefixKVs, keyvals...)
	slsPrefixKVs = append(slsPrefixKVs, keyvals...)
}

func ExtractCtx(ctx context.Context, logType LogType) (keyvals []interface{}) {
	prefixs := zapPrefixKVs
	if logType == LogTypeSls {
		prefixs = slsPrefixKVs
	}
	kvs := make([]interface{}, 0, len(prefixs))
	kvs = append(kvs, prefixs...)
	for i := 1; i < len(kvs); i += 2 {
		if v, ok := kvs[i].(Valuer); ok {
			kvs[i] = v(ctx)
		}
	}
	return kvs
}

// Timestamp returns a timestamp Valuer with a custom time format.
func Timestamp(layout string) Valuer {
	return func(context.Context) interface{} {
		return time.Now().Format(layout)
	}
}

// Caller returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
		idx = strings.LastIndexByte(file[:idx], '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

func TraceID() Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			return span.TraceID().String()
		}
		return ""
	}
}
