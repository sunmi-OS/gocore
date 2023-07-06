package http_request

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sunmi-OS/gocore/v2/utils"
)

type HttpClient struct {
	Client  *resty.Client
	Request *resty.Request

	enableMessageCodeMetrics bool // default: false
	disableMetrics           bool // default: false
	disableLog               bool // default: false
	disableBreaker           bool // default: true
}

type option struct {
	slowThresholdMs          int64
	hideRespBodyLogsWithPath map[string]bool
	hideReqBodyLogsWithPath  map[string]bool
}

type Option func(op *option)

// New 重试3次、间隔2~10s、非2xx状态码抛异常到error里、默认打开统计、日志、bbr限流(待实现)、统计业务错误码
func New(options ...Option) *HttpClient {
	// Create a Resty Client
	client := resty.New()

	// Retries are configured per client
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// TimeOut
		SetTimeout(5*time.Second).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(2*time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(5*time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		}).
		OnAfterResponse(MustCode200).
		SetHeader(utils.XAppName, utils.GetAppName())

	c := &HttpClient{
		Client:                   client,
		disableMetrics:           false,
		disableLog:               false,
		enableMessageCodeMetrics: false,
		disableBreaker:           true, // default disable, will open soon
	}
	c.setLog(options...)
	return c
}

func (h *HttpClient) SetTrace(header interface{}) *HttpClient {
	trace := utils.SetHeader(header)
	h.Request.Header = trace.HttpHeader
	return h
}

func (h *HttpClient) SetDisableMetrics(disable bool) *HttpClient {
	h.disableMetrics = disable
	return h
}

func (h *HttpClient) SetDisableLog(disable bool) *HttpClient {
	h.disableLog = disable
	return h
}

func (h *HttpClient) SetDisableBreaker(disable bool) *HttpClient {
	h.disableBreaker = disable
	return h
}

func (h *HttpClient) SetEnableMessageCodeMetrics(enable bool) *HttpClient {
	h.enableMessageCodeMetrics = enable
	return h
}

// ErrIncorrectCode 非2xx 状态码
var ErrIncorrectCode = errors.New("incorrect http status")

// MustCode200 将非200 状态码认为错误的 middleware (注意：设置 SetDoNotParseResponse 后不会触发任何 midware，包括这个)
func MustCode200(cli *resty.Client, resp *resty.Response) error {
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%w:%d host:%s, url:%s", ErrIncorrectCode, resp.StatusCode(),
			resp.Request.RawRequest.URL.Host, resp.Request.RawRequest.URL.Path,
		)
	}
	return nil
}
