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

func (h *HttpClient) setLog(options ...Option) *HttpClient {
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

		fields := []interface{}{
			"logtype", "http_client",
			"cost", resp.Time().Milliseconds(),
			"traceid", resp.Header().Get(utils.XB3TraceId),

			"r_method", r.Method,
			"r_host", r.RawRequest.URL.Host,
			"r_path", path,
			"r_body", reqBody,
			"r_header", r.Header,
			"r_start", r.Time.Format(utils.TimeFormat),

			"s_body", respBody,
			"s_status", resp.StatusCode(),
			"s_header", resp.Header(),
		}
		param := r.QueryParam.Encode()
		if param != "" {
			fields = append(fields, "r_params", param)
		}
		_ = r.RawRequest.RemoteAddr

		clientSendBytes.WithLabelValues(path).Add(mustPositive(sendBytes))
		clientRecvBytes.WithLabelValues(path).Add(mustPositive(recvBytes))
		if !h.disableMetrics {
			clientReqDur.WithLabelValues(path).Observe(float64(time.Since(r.Time).Milliseconds()))
			statusCode := resp.StatusCode()
			clientReqCodeTotal.WithLabelValues(path, strconv.FormatInt(int64(statusCode), 10)).Inc()
			if statusCode == http.StatusOK && h.enableMessageCodeMetrics {
				if root, err0 := sonic.Get(resp.Body()); err0 == nil {
					respCode, _ := root.Get("code").Int64()
					respMsg, _ := root.Get("msg").String()
					fields = append(fields, "s_code", respCode)
					fields = append(fields, "s_msg", respMsg)
				}
			}
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
