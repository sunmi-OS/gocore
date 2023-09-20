package http_otel

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type OtelRoundTripper struct {
	original http.RoundTripper
}

func (ort *OtelRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(req.Context(), req.URL.Path)
	defer span.End()

	req = req.WithContext(ctx)
	resp, err := ort.original.RoundTrip(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		attrs := []attribute.KeyValue{
			{
				Key:   "http.status_code",
				Value: attribute.IntValue(resp.StatusCode),
			},
			{
				Key:   "http.method",
				Value: attribute.StringValue(req.Method),
			},
			{
				Key:   "http.route",
				Value: attribute.StringValue(req.URL.RequestURI()),
			},
			{
				Key:   "http.scheme",
				Value: attribute.StringValue(req.URL.Scheme),
			},
			{
				Key:   "http.user_agent",
				Value: attribute.StringValue(req.UserAgent()),
			},
		}
		span.SetAttributes(attrs...)
	}
	return resp, err
}

func New(tripper http.RoundTripper) *OtelRoundTripper {
	return &OtelRoundTripper{
		original: tripper,
	}
}
