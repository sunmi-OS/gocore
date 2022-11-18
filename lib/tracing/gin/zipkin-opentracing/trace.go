package zipkin_opentracing

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	zipkinopentracing "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

// ZipKinOpentracing zipkin+openTracing
// serviceName 接入tracing 的服务名称
// exportType 上报类型 1-实时上报 2-先转换为本地日志文件，暂时只支持实时上报
// endPointUrl 链路日志上报地址,exportType为2时无任何作用
func ZipKinOpentracing(serviceName string, exportType int64, endPointUrl string) gin.HandlerFunc {
	appName := serviceName + "#" + utils.GetRunTime()
	tracer := newZipKinTracer(appName, endPointUrl)
	mwOptions := ginhttp.OperationNameFunc(func(r *http.Request) string {
		return r.URL.Path
	})
	spanObserver := ginhttp.MWSpanObserver(func(span opentracing.Span, r *http.Request) {
		span.SetTag("monitor.api", serviceName)
	})
	return ginhttp.Middleware(tracer, mwOptions, spanObserver)
}

// newZipKinTracer 通过 http 直接上报
// appName application name
func newZipKinTracer(appName string, endPointUrl string) opentracing.Tracer {
	//zipkinreporter.Timeout 上报链路日志超时时间（http）
	//zipkinreporter.BatchSize 每次推送数量
	//zipkinreporter.BatchInterval 批量推送周期
	//zipkinreporter.MaxBacklog 链路日志缓冲区大小，最大1000，超过1000会被丢弃
	reporter := zipkinreporter.NewReporter(endPointUrl, zipkinreporter.Timeout(2*time.Second))
	//create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(appName, "localhost:0")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinopentracing.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(tracer)
	return tracer
}

// StartSpanWithCtx 生成上下文span
// skip The argument skip is the number of stack frames to ascend, with 0 identifying the caller of Caller
func StartSpanWithCtx(ctx context.Context, skip int) (opentracing.Span, context.Context) {
	//开始链路追踪
	pc, _, _, _ := runtime.Caller(skip)
	spanName := ""
	if pc > 0 {
		spanName = spanName + "/" + runtime.FuncForPC(pc).Name()
	}
	return opentracing.StartSpanFromContext(ctx, spanName)
}
