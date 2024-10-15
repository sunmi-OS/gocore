package logx

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/sunmi-OS/gocore/v2/api"
	"github.com/sunmi-OS/gocore/v2/utils"
)

type Valuer func(ctx context.Context) interface{}

var zapPrefixKVs = []interface{}{
	"traceid", GetCtxKey(utils.XB3TraceId),
	"appname", utils.GetAppName(),
	"runtime", utils.GetRunTime(),
	"hostname", utils.GetHostname(),
}

var slsPrefixKVs = []interface{}{
	"ts", Timestamp(utils.TimeFormat),
	"caller", Caller(6),
	"traceid", GetCtxKey(utils.XB3TraceId),
	"appname", utils.GetAppName(),
	"runtime", utils.GetRunTime(),
	"hostname", utils.GetHostname(),
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

func SetCtxKV(ctx context.Context, key, val string) context.Context {
	return SetCtxKVS(ctx, map[string]string{key: val})
}

func SetCtxKVS(ctx context.Context, kvs map[string]string) context.Context {
	if apiCtx, ok := ctx.(*api.Context); ok {
		apiCtx.Request = apiCtx.Request.WithContext(utils.SetMetaDataMulti(apiCtx.Request.Context(), kvs))
		return apiCtx.Request.Context()
	}
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Request = ginCtx.Request.WithContext(utils.SetMetaDataMulti(ginCtx.Request.Context(), kvs))
		return ginCtx.Request.Context()
	}
	return utils.SetMetaDataMulti(ctx, kvs)
}

func GetCtxKey(key string) Valuer {
	return func(ctx context.Context) interface{} {
		switch c := ctx.(type) {
		case api.Context:
			return utils.GetMetaData(c.Request.Context(), key)
		case *api.Context:
			return utils.GetMetaData(c.Request.Context(), key)
		case *gin.Context:
			return utils.GetMetaData(c.Request.Context(), key)
		default:
			return utils.GetMetaData(ctx, key)
		}
	}
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
