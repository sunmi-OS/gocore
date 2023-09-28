package zipkin_opentracing

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type RoundTripper struct {
	original http.RoundTripper
}

func (ort *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	span, ctx := opentracing.StartSpanFromContext(req.Context(), req.Method+" "+req.URL.Path)
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	defer span.Finish()
	req = req.WithContext(ctx)
	resp, err := ort.original.RoundTrip(req)
	if err != nil {
		span.SetTag("error", err.Error())
		span.LogFields(
			log.String("event", "error"),
			log.String("message", "Something went wrong"),
			log.Error(err),
		)
	} else {
		span.SetTag("http.status_code", resp.StatusCode)
		span.SetTag("http.method", req.Method)
		span.SetTag("http.route", req.URL.RequestURI())
		span.SetTag("http.scheme", req.URL.Scheme)
		span.SetTag("http.user_agent", req.UserAgent())
	}
	return resp, err
}

func NewRoundTripper(tripper http.RoundTripper) *RoundTripper {
	return &RoundTripper{
		original: tripper,
	}
}
