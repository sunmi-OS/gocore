package trace

import (
	"github.com/uber/jaeger-client-go"
)

const (
	// Trace.Id
	TagTraceId = "Trace-Id"

	// The software package, framework, library, or module that generated the associated Span.
	// E.g., "grpc", "django", "JDBI".
	// type string
	TagComponent = "component"

	// HTTP method of the request for the associated Span. E.g., "GET", "POST"
	// type string
	TagHTTPMethod = "http.method"

	// URL of the request being handled in this segment of the trace, in standard URI format.
	// E.g., "https://domain.net/path/to?resource=here"
	// type string
	TagHTTPURL = "http.url"

	TagHTTPRaw = "http.raw"

	TagHTTPBody = "http.body"
)

type Config struct {
	ServiceName string
	Endpoint    string
	TraceOpts   []jaeger.TracerOption
}
