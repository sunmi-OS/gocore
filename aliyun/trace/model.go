package trace

import (
	"github.com/uber/jaeger-client-go"
)

const (
	TagTraceId = "Trace-Id"

	TagComponent = "component"

	TagHTTPMethod = "http.method"

	TagHTTPURL = "http.url"

	TagHTTPRaw = "http.raw"

	TagHTTPBody = "http.body"
)

type Config struct {
	ServiceName string
	Endpoint    string
	TraceOpts   []jaeger.TracerOption
}
