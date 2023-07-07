package accesslog

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const hideBody = "gocore_body_hide"

// ServerLogging middleware for accesslog
func ServerLogging(options ...Option) gin.HandlerFunc {
	op := option{
		slowThresholdMs:     1000,
		hideLogsWithPath:    hideLogsPath,
		hideReqBodyWithPath: nil,
		hideRespBodWithPath: hidelRespBodyLogsPath,
		allowShowHeaders:    map[string]bool{},
		hideShowHeaders:     hideShowHeaders,
	}
	for _, apply := range options {
		apply(&op)
	}

	return func(c *gin.Context) {
		r := c.Request
		path := r.URL.Path
		start := time.Now()
		quota := int64(-1)
		if deadline, ok := r.Context().Deadline(); ok {
			quota = time.Until(deadline).Milliseconds()
		}
		body := ""
		if !op.hideReqBodyWithPath[path] {
			b, err := c.GetRawData()
			if err != nil {
				body = "failed to get request body"
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				body = string(b)
			}
		} else {
			body = hideBody
		}

		hideResp := op.hideRespBodWithPath[path]
		var writer responseWriter
		if !hideResp {
			writer = responseWriter{
				c.Writer,
				bytes.NewBuffer([]byte{}),
			}
			c.Writer = writer
		}

		c.Next()

		r = c.Request
		ctx := r.Context()
		responseCode := math.MinInt8
		var responseMsg string
		var respBytes []byte
		if !hideResp {
			respBytes = writer.b.Bytes()
			if root, err0 := sonic.Get(respBytes); err0 == nil {
				code, err := root.Get("code").Int64()
				if err == nil {
					responseCode = int(code)
				}
				msg, err := root.Get("msg").String()
				if err == nil {
					responseMsg = msg
				}
			}
		} else {
			respBytes = []byte(hideBody)
		}

		sendBytes := mustPositive(float64(c.Writer.Size()))
		recvBytes := mustPositive(float64(r.ContentLength))
		reqAppname := r.Header.Get(utils.XAppName)
		statusCode := c.Writer.Status()
		serverRecvBytes.WithLabelValues(path, reqAppname).Add(recvBytes)
		serverSendBytes.WithLabelValues(path, reqAppname).Add(sendBytes)
		serverReqCodeTotal.WithLabelValues(path, reqAppname, strconv.FormatInt(int64(statusCode), 10)).Inc()
		costms := time.Since(start).Milliseconds()
		serverReqDur.WithLabelValues(path, reqAppname).Observe(float64(costms))

		if op.hideLogsWithPath[path] {
			return
		}

		fields := []interface{}{
			"kind", "server",
			"costms", costms,
			"traceid", c.GetHeader(utils.XB3TraceId), // 后续待优化
			"ip", c.ClientIP(),
			"host", r.Host,
			"method", r.Method,
			"path", path,
			"req", body,
			"resp", string(respBytes),
			"status", statusCode, // http状态码
			"code", responseCode, // 业务错误码
			"msg", responseMsg,

			"start_time", start.Format(utils.TimeFormat),
			"req_header", filterHeaders(r.Header, op.allowShowHeaders, op.hideShowHeaders),
			"resp_header", c.Writer.Header(),
		}
		if reqAppname != "" {
			fields = append(fields, "req_appname", reqAppname)
		}
		if r.URL.RawQuery != "" {
			fields = append(fields, "params", r.URL.RawQuery)
		}
		if c.GetHeader("x-forwarded-for") != "" {
			fields = append(fields, "forward_ip", c.GetHeader("x-forwarded-for"))
		}
		if quota != -1 {
			fields = append(fields, "timeout_quota", quota) // 收到请求时，剩余处理时间
		}

		logFunc := glog.InfoV
		if statusCode >= http.StatusInternalServerError {
			logFunc = glog.ErrorV
		} else if statusCode >= http.StatusBadRequest {
			logFunc = glog.WarnV
		} else if op.slowThresholdMs != 0 && costms > op.slowThresholdMs {
			logFunc = glog.WarnV
		}
		logFunc(ctx, fields...)
	}
}

func mustPositive(val float64) float64 {
	if val < 0 {
		return 0
	}
	return val
}

// header白名单过滤
func filterHeaders(headers http.Header, allowShowHeaders map[string]bool, hideShowHeaders map[string]bool) http.Header {
	filteredHeaders := http.Header{}
	for k, v := range headers {
		lower := strings.ToLower(k)
		// 优先判断白名单
		if allowShowHeaders[k] {
			filteredHeaders[k] = v
			continue
		}
		if hideShowHeaders[lower] {
			continue
		}
		// 如果没配置白名单，能走到这里则允许打印
		if len(allowShowHeaders) == 0 {
			filteredHeaders[k] = v
		}
	}
	return filteredHeaders
}

const (
	serverNamespace = "http_server"
)

var (
	serverReqDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
	}, []string{"path", "caller"})
	serverReqCodeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
	}, []string{"path", "caller", "code"})
	serverSendBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serverNamespace,
		Subsystem: "bandwith",
		Name:      "send",
	}, []string{"method", "caller"})
	serverRecvBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serverNamespace,
		Subsystem: "bandwith",
		Name:      "recv",
	}, []string{"method", "caller"})
)

func init() {
	prometheus.MustRegister(serverReqDur)
	prometheus.MustRegister(serverReqCodeTotal)
	prometheus.MustRegister(serverSendBytes)
	prometheus.MustRegister(serverRecvBytes)
}

// 黑名单 某些路径不打印response body，但打印日志
var hidelRespBodyLogsPath = map[string]bool{
	"/debug/pprof/":        true,
	"/debug/pprof/cmdline": true,
	"/debug/pprof/profile": true,
	"/debug/pprof/symbol":  true,
	"/debug/pprof/trace":   true,
}

// 黑名单 某些路径不打印日志
var hideLogsPath = map[string]bool{
	"/metrics": true,
	"/health":  true,
}

var hideShowHeaders = map[string]bool{
	"accept":          true,
	"accept-encoding": true,
}

// WithSlowThreshold 当请求耗时超过slowThreshold时，打印slow log。建议配置1000
func WithSlowThreshold(slowThresholdMs int64) Option {
	return func(o *option) {
		o.slowThresholdMs = slowThresholdMs
	}
}

// WithHideLogsPath 对某些路径不打印日志
func WithHideLogsPath(hideLogsWithPath map[string]bool, isMerge bool) Option {
	return func(o *option) {
		if isMerge {
			o.hideLogsWithPath = mergeMap(o.hideLogsWithPath, hideLogsWithPath)
		} else {
			o.hideLogsWithPath = hideLogsWithPath
		}
	}
}

// WithHideBodyLogsPath 对某些路径不打印body
func WithHideBodyLogsPath(hideBodyLogsWithPath map[string]bool, isMerge bool) Option {
	return func(o *option) {
		if isMerge {
			o.hideRespBodWithPath = mergeMap(o.hideRespBodWithPath, hideBodyLogsWithPath)
		} else {
			o.hideRespBodWithPath = hideBodyLogsWithPath
		}
	}
}

// WithAllowShowHeaders 只展示某些header
func WithAllowShowHeaders(allowHeaders []string) Option {
	return func(o *option) {
		for _, header := range allowHeaders {
			o.allowShowHeaders[strings.ToLower(header)] = true
		}
	}
}

func WithHideShowHeaders(hideHeaders map[string]bool, isMerge bool) Option {
	return func(o *option) {
		if isMerge {
			o.hideShowHeaders = mergeMap(o.hideShowHeaders, hideHeaders)
		} else {
			o.hideShowHeaders = hideHeaders
		}
	}
}

func mergeMap(m1, m2 map[string]bool) map[string]bool {
	for k, v := range m2 {
		if _, ok := m1[k]; !ok {
			m1[k] = v
		}
	}
	return m1
}

type option struct {
	slowThresholdMs     int64
	hideLogsWithPath    map[string]bool
	hideReqBodyWithPath map[string]bool
	hideRespBodWithPath map[string]bool
	allowShowHeaders    map[string]bool
	hideShowHeaders     map[string]bool
}

type Option func(op *option)
