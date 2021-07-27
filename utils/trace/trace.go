package trace

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

type Tracer struct {
	ServiceName string
	Trace       opentracing.Tracer
	closer      io.Closer
}

func NewTracer(c *Config) (trace *Tracer) {
	ht := transport.NewHTTPTransport(c.Endpoint)
	tracer, closer := jaeger.NewTracer(
		c.ServiceName,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(ht),
		c.TraceOpts...,
	)
	trace = &Tracer{
		ServiceName: c.ServiceName,
		Trace:       tracer,
		closer:      closer,
	}
	return
}

func (t *Tracer) GinTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// span
		spanCtx, _ := t.Trace.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		span := t.Trace.StartSpan(c.FullPath(), opentracing.ChildOf(spanCtx))
		defer span.Finish()
		traceId := uuid.New().String()
		span.SetTag(TagTraceId, traceId)
		span.SetTag(TagComponent, "net/http")
		span.SetTag(TagHTTPMethod, c.Request.Method)
		span.SetTag(TagHTTPURL, c.Request.URL.String())
		rawReq, _ := httputil.DumpRequest(c.Request, false)
		span.SetTag(TagHTTPRaw, string(rawReq))
		switch c.Request.Method {
		case http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodPut:
			bs, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
			span.SetTag(TagHTTPBody, string(bs))
		default:
		}
		// export trace id to user.
		c.Writer.Header().Set(TagTraceId, traceId)
		c.Next()
	}
}

func (t *Tracer) Close() {
	t.closer.Close()
}
