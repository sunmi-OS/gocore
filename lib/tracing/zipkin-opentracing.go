package tracing

import (
	"context"
	zipkinopentracing "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sunmi-OS/gocore/utils"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const defaultComponentName = "net/http"

type mwOptions struct {
	opNameFunc    func(r *http.Request) string
	spanObserver  func(span opentracing.Span, r *http.Request)
	urlTagFunc    func(u *url.URL) string
	componentName string
}

// mwOption controls the behavior of the Middleware.
type mwOption func(*mwOptions)

// operationNameFunc returns a mwOption that uses given function f
// to generate operation name for each server-side span.
func operationNameFunc(f func(r *http.Request) string) mwOption {
	return func(options *mwOptions) {
		options.opNameFunc = f
	}
}

// mwSpanObserver returns a mwOption that observe the span
// for the server-side span.
func mwSpanObserver(f func(span opentracing.Span, r *http.Request)) mwOption {
	return func(options *mwOptions) {
		options.spanObserver = f
	}
}

func middleware(tr opentracing.Tracer, options ...mwOption) echo.MiddlewareFunc {
	opts := mwOptions{
		opNameFunc: func(r *http.Request) string {
			return "HTTP " + r.Method
		},
		spanObserver: func(span opentracing.Span, r *http.Request) {},
		urlTagFunc: func(u *url.URL) string {
			return u.String()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			carrier := opentracing.HTTPHeadersCarrier(c.Request().Header)
			ctx, _ := tr.Extract(opentracing.HTTPHeaders, carrier)
			op := opts.opNameFunc(c.Request())
			sp := tr.StartSpan(op, ext.RPCServerOption(ctx))
			ext.HTTPMethod.Set(sp, c.Request().Method)
			ext.HTTPUrl.Set(sp, opts.urlTagFunc(c.Request().URL))
			opts.spanObserver(sp, c.Request())

			// set component name, use "net/http" if caller does not specify
			componentName := opts.componentName
			if componentName == "" {
				componentName = defaultComponentName
			}
			ext.Component.Set(sp, componentName)

			c.SetRequest(c.Request().WithContext(
				opentracing.ContextWithSpan(c.Request().Context(), sp)))

			err := next(c)
			if err != nil {
				return err
			}
			ext.HTTPStatusCode.Set(sp, uint16(c.Response().Status))
			sp.Finish()
			return nil
		}
	}
}

// 通过 SDK 直接上报 (HTTP)
func newZipKinTracer(serviceName string, endPointUrl string) opentracing.Tracer {
	//zipkinreporter.Timeout 上报链路日志超时时间（http）
	//zipkinreporter.BatchSize 每次推送数量
	//zipkinreporter.BatchInterval 批量推送周期
	//zipkinreporter.MaxBacklog 链路日志缓冲区大小，最大1000，超过1000会被丢弃
	reporter := zipkinreporter.NewReporter(endPointUrl, zipkinreporter.Timeout(1*time.Second))

	//create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, "localhost:0")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	//zipkin.NewModuloSampler
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

func ZipKinOpentracing(serviceName string, endPointUrl string) echo.MiddlewareFunc {
	appName := serviceName + "#" + utils.GetRunTime()
	tracer := newZipKinTracer(appName, endPointUrl)
	mwOptions := operationNameFunc(func(r *http.Request) string {
		return r.URL.Path
	})
	//监控报警
	spanObserver := mwSpanObserver(func(span opentracing.Span, r *http.Request) {
		span.SetTag("monitor.api", serviceName)
	})
	return middleware(tracer, mwOptions, spanObserver)
}

// StartSpanWithCtx 生成上下文span
// skip The argument skip is the number of stack frames to ascend, with 0 identifying the caller of Caller
func StartSpanWithCtx(ctx context.Context, skip int) (opentracing.Span, context.Context) {
	pc, _, _, _ := runtime.Caller(skip)
	spanName := ""
	if pc > 0 {
		spanName = spanName + "/" + runtime.FuncForPC(pc).Name()
	}
	return opentracing.StartSpanFromContext(ctx, spanName)
}
