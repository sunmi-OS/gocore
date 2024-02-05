package http_request

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sunmi-OS/gocore/v2/utils"
)

const maxShowBodySize = 1024 * 100

type HttpClient struct {
	Client  *resty.Client
	Request *resty.Request

	disableLog               bool  // default: false 默认打印日志(配置SetLog后)
	disableMetrics           bool  // default: false 默认开启统计
	disableBreaker           bool  // default: true 默认关闭熔断
	slowThresholdMs          int64 // default: 0 默认关闭慢请求打印
	hideRespBodyLogsWithPath map[string]bool
	hideReqBodyLogsWithPath  map[string]bool
	maxShowBodySize          int64
}

func New() *HttpClient {
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
		SetHeader(utils.XAppName, utils.GetAppName())

	return &HttpClient{
		Client:                   client,
		Request:                  client.R(),
		disableMetrics:           false,
		disableLog:               false,
		disableBreaker:           true, // default disable, will open soon
		hideReqBodyLogsWithPath:  hidelBodyLogsPath,
		hideRespBodyLogsWithPath: hidelBodyLogsPath,
		maxShowBodySize:          maxShowBodySize,
	}
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

func (h *HttpClient) SetMaxShowBodySize(bodySize int64) *HttpClient {
	h.maxShowBodySize = bodySize
	return h
}

func (h *HttpClient) SetRespBodyLogsWithPath(paths []string) *HttpClient {
	h.hideRespBodyLogsWithPath = make(map[string]bool)
	for _, path := range paths {
		h.hideRespBodyLogsWithPath[path] = true
	}
	return h
}

func (h *HttpClient) SetReqBodyLogsWithPath(paths []string) *HttpClient {
	h.hideReqBodyLogsWithPath = make(map[string]bool)
	for _, path := range paths {
		h.hideReqBodyLogsWithPath[path] = true
	}
	return h
}

func (h *HttpClient) SetSlowThresholdMs(threshold int64) *HttpClient {
	h.slowThresholdMs = threshold
	return h
}

// ErrIncorrectCode 非2xx 状态码
var ErrIncorrectCode = errors.New("incorrect http status")

// MustCode200 将非200的状态码认为错误的(注意：设置 SetDoNotParseResponse 后不会触发任何middware包括这个)
func MustCode200(cli *resty.Client, resp *resty.Response) error {
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%w:%d host:%s, url:%s", ErrIncorrectCode, resp.StatusCode(),
			resp.Request.RawRequest.URL.Host, resp.Request.RawRequest.URL.Path,
		)
	}
	return nil
}

// 黑名单 某些路径不打印response body，但打印日志
var hidelBodyLogsPath = map[string]bool{
	"/debug/pprof/":        true,
	"/debug/pprof/cmdline": true,
	"/debug/pprof/profile": true,
	"/debug/pprof/symbol":  true,
	"/debug/pprof/trace":   true,
}
