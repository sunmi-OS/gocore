package http_request

import (
	"errors"
	"time"

	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/go-resty/resty/v2"
)

type HttpClient struct {
	Client  *resty.Client
	Request *resty.Request
}

func New() *HttpClient {

	// Create a Resty Client
	client := resty.New()

	// Retries are configured per client
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(10).
		// TimeOut
		SetTimeout(5 * time.Second).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(2 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20 * time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})

	return &HttpClient{
		Client:  client,
		Request: client.R(),
	}
}

func (h *HttpClient) SetTrace(header interface{}) *HttpClient {
	trace := utils.SetHeader(header)

	h.Request.SetHeader(utils.XRequestId, trace.HttpHeader.Get(utils.XRequestId))
	h.Request.SetHeader(utils.XB3TraceId, trace.HttpHeader.Get(utils.XB3TraceId))
	h.Request.SetHeader(utils.XB3SpanId, trace.HttpHeader.Get(utils.XB3SpanId))
	h.Request.SetHeader(utils.XB3ParentSpanId, trace.HttpHeader.Get(utils.XB3ParentSpanId))
	h.Request.SetHeader(utils.XB3Sampled, trace.HttpHeader.Get(utils.XB3Sampled))
	h.Request.SetHeader(utils.XB3Flags, trace.HttpHeader.Get(utils.XB3Flags))
	h.Request.SetHeader(utils.B3, trace.HttpHeader.Get(utils.B3))
	h.Request.SetHeader(utils.XOtSpanContext, trace.HttpHeader.Get(utils.XOtSpanContext))

	h.Request.Header = trace.HttpHeader
	return h
}
