package http_request

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/glog/logx"
	"github.com/sunmi-OS/gocore/v2/glog/sls"
	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

const hideBody = "gocore_body_hide"

type Log interface {
	InfoV(ctx context.Context, keyvals ...interface{}) error
	WarnV(ctx context.Context, keyvals ...interface{}) error
	ErrorV(ctx context.Context, keyvals ...interface{}) error
}

func (h *HttpClient) SetLog(log Log) *HttpClient {
	h.Client = h.Client.OnBeforeRequest(func(client *resty.Client, r *resty.Request) error {
		traceid := utils.GetMetaData(r.Context(), utils.XB3TraceId)
		if traceid != "" {
			r = r.SetHeader(utils.XB3TraceId, traceid)
		}
		return nil
	})

	h.Client = h.Client.OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
		r := resp.Request
		ctx := resp.Request.Context()
		ctx = utils.SetMetaData(ctx, utils.XB3TraceId, resp.Header().Get(utils.XB3TraceId))
		var reqBody interface{}
		reqBody = hideBody
		respBody := hideBody
		path := r.RawRequest.URL.Path
		sendBytes := r.RawRequest.ContentLength
		recvBytes := resp.Size()
		statusCode := resp.StatusCode()
		if !h.hideRespBodyLogsWithPath[path] && recvBytes < h.maxShowBodySize {
			respBody = string(resp.Body())
		}
		if !h.hideReqBodyLogsWithPath[path] && sendBytes < h.maxShowBodySize {
			reqBody = r.Body
		}

		if !h.disableMetrics {
			clientSendBytes.WithLabelValues(path).Add(mustPositive(sendBytes))
			clientRecvBytes.WithLabelValues(path).Add(mustPositive(recvBytes))
			clientReqDur.WithLabelValues(path).Observe(float64(time.Since(r.Time).Milliseconds()))
			clientReqCodeTotal.WithLabelValues(path, strconv.FormatInt(int64(statusCode), 10)).Inc()
		}

		if h.disableLog {
			return nil
		}
		fields := []interface{}{
			"kind", "client",
			"costms", resp.Time().Milliseconds(),
			"method", r.Method,
			"host", r.RawRequest.URL.Host,
			"path", path,
			"req", reqBody,
			"resp", utils.LogContentUnmarshal(respBody),
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
		fields = append(fields, "start_time", r.Time.Format(utils.TimeFormat))
		fields = append(fields, "req_header", r.Header)
		fields = append(fields, "resp_header", resp.Header())
		param := r.QueryParam.Encode()
		if param != "" {
			fields = append(fields, "params", param)
		}

		logFunc := log.InfoV
		if resp.StatusCode() >= http.StatusInternalServerError {
			logFunc = log.ErrorV
		} else if resp.StatusCode() >= http.StatusBadRequest {
			logFunc = log.WarnV
		}
		logFunc(ctx, fields...)
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

type GocoreLog struct {
}

func NewGocoreLog() *GocoreLog {
	return &GocoreLog{}
}

func (l *GocoreLog) InfoV(ctx context.Context, keyvals ...interface{}) error {
	glog.InfoV(ctx, keyvals...)
	return nil
}
func (l *GocoreLog) WarnV(ctx context.Context, keyvals ...interface{}) error {
	glog.WarnV(ctx, keyvals...)
	return nil
}
func (l *GocoreLog) ErrorV(ctx context.Context, keyvals ...interface{}) error {
	glog.ErrorV(ctx, keyvals...)
	return nil
}

func NewAliyunLog(topic string) *AliyunLog {
	return &AliyunLog{topic: topic}
}

type AliyunLog struct {
	topic string
}

func (l *AliyunLog) InfoV(ctx context.Context, keyvals ...interface{}) error {
	ctx = utils.SetMetaData(ctx, logx.SlsTopic, l.topic)
	return sls.LogClient.CommonLog(logx.LevelInfo, ctx, keyvals...)
}

func (l *AliyunLog) WarnV(ctx context.Context, keyvals ...interface{}) error {
	ctx = utils.SetMetaData(ctx, logx.SlsTopic, l.topic)
	return sls.LogClient.CommonLog(logx.LevelWarn, ctx, keyvals...)
}

func (l *AliyunLog) ErrorV(ctx context.Context, keyvals ...interface{}) error {
	ctx = utils.SetMetaData(ctx, logx.SlsTopic, l.topic)
	return sls.LogClient.CommonLog(logx.LevelError, ctx, keyvals...)
}
