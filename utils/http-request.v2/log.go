package http_request

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

// 黑名单 某些路径不打印response body，但打印日志
var hidelBodyLogsPath = map[string]bool{
	"/debug/pprof/":        true,
	"/debug/pprof/cmdline": true,
	"/debug/pprof/profile": true,
	"/debug/pprof/symbol":  true,
	"/debug/pprof/trace":   true,
}

const hideBody = "gocore_body_hide"

func (h *HttpClient) SetLog(options ...Option) *HttpClient {
	op := option{
		slowThresholdMs:          1000,
		hideRespBodyLogsWithPath: hidelBodyLogsPath,
	}
	for _, apply := range options {
		apply(&op)
	}

	h.Client = h.Client.OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
		r := resp.Request
		ctx := resp.Request.Context()
		var reqBody interface{}
		reqBody = hideBody
		respBody := hideBody
		path := r.RawRequest.URL.Path
		if !op.hideRespBodyLogsWithPath[path] {
			respBody = string(resp.Body())
		}
		if !op.hideReqBodyLogsWithPath[path] {
			reqBody = r.Body
		}
		sendBytes := r.RawRequest.ContentLength
		recvBytes := resp.Size()
		statusCode := resp.StatusCode()

		fields := []interface{}{
			"kind", "client",
			"costms", resp.Time().Milliseconds(),
			"traceid", resp.Header().Get(utils.XB3TraceId),
			"method", r.Method,
			"host", r.RawRequest.URL.Host,
			"path", path,
			"req", reqBody,
			"resp", respBody,
			"status", statusCode,
		}
		if statusCode == http.StatusOK {
			if root, err0 := sonic.Get(resp.Body()); err0 == nil {
				respCode, err := root.Get("code").Int64()
				if err == nil {
					fields = append(fields, "code", respCode)
				}
				respMsg, err := root.Get("msg").String()
				if err == nil {
					fields = append(fields, "msg", respMsg)
				}
			}
		}
		clientSendBytes.WithLabelValues(path).Add(mustPositive(sendBytes))
		clientRecvBytes.WithLabelValues(path).Add(mustPositive(recvBytes))
		if !h.disableMetrics {
			clientReqDur.WithLabelValues(path).Observe(float64(time.Since(r.Time).Milliseconds()))
			clientReqCodeTotal.WithLabelValues(path, strconv.FormatInt(int64(statusCode), 10)).Inc()

		}
		fields = append(fields, "start_time", r.Time.Format(utils.TimeFormat))
		fields = append(fields, "req_header", r.Header)
		fields = append(fields, "resp_header", resp.Header())
		param := r.QueryParam.Encode()
		if param != "" {
			fields = append(fields, "params", param)
		}

		if !h.disableLog {
			logFunc := glog.InfoV
			if resp.StatusCode() >= http.StatusInternalServerError {
				logFunc = glog.ErrorV
			} else if resp.StatusCode() >= http.StatusBadRequest {
				logFunc = glog.WarnV
			}
			logFunc(ctx, fields...)
		}
		return nil
	})
	return h
}

func mustPositive(val int64) float64 {
	if val < 0 {
		return 0
	}
	return float64(val)
}
