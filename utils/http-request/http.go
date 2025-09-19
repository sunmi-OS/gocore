package http_request

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sunmi-OS/gocore/v2/utils"
	"google.golang.org/grpc/metadata"
)

const maxShowBodySize = 1024 * 100

type HttpClient struct {
	Client                   *resty.Client
	disableLog               bool            // default: false 默认打印日志(配置SetLog后)
	disableMetrics           bool            // default: false 默认开启统计
	disableBreaker           bool            // default: true 默认关闭熔断
	slowThresholdMs          int64           // default: 0 默认关闭慢请求打印
	hideRespBodyLogsWithPath map[string]bool // 不打印path在map里的返回体
	hideReqBodyLogsWithPath  map[string]bool // 不打印path在map里的请求体
	maxShowBodySize          int64
}

// Option support set var to resty.Client
type Option func(client *resty.Client)

func New(opts ...Option) *HttpClient {
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

	for _, opt := range opts {
		opt(client)
	}

	return &HttpClient{
		Client:                   client,
		disableMetrics:           false,
		disableLog:               false,
		disableBreaker:           true, // default disable, will open soon
		hideReqBodyLogsWithPath:  hidelBodyLogsPath,
		hideRespBodyLogsWithPath: hidelBodyLogsPath,
		maxShowBodySize:          maxShowBodySize,
	}
}

// R call Client.R() , and set the context value to headers
func (h *HttpClient) R(ctx context.Context) *resty.Request {
	req := h.Client.R()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// if mdIncomingKey is absent, return resty Client
		return req
	}

	req.SetHeaderMultiValues(md)

	// when call downstream service, set current app name as client app
	req.SetHeader(utils.XClientApp, utils.GetAppName())

	return req
}

// WithRetryWaitTime set retry wait time
func WithRetryWaitTime(retry int) Option {
	return func(hc *resty.Client) {
		hc.SetRetryWaitTime(time.Duration(retry) * time.Second)
	}
}

// WithRetryMaxWaitTime set retry max wait time
func WithRetryMaxWaitTime(retry int) Option {
	return func(hc *resty.Client) {
		hc.SetRetryMaxWaitTime(time.Duration(retry) * time.Second)
	}
}

// WithRetryCount set retry times
func WithRetryCount(retry int) Option {
	return func(hc *resty.Client) {
		hc.SetRetryCount(retry)
	}
}

// WithTimeout set request timeout, unit is second
func WithTimeout(timeout int64) Option {
	return func(hc *resty.Client) {
		hc.SetTimeout(time.Duration(timeout) * time.Second)
	}
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
