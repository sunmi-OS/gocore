package web

import (
	"github.com/sunmi-OS/gocore/v2/utils/limit"
	"github.com/sunmi-OS/gocore/v2/utils/trace"
)

type Config struct {
	// http host
	Host string
	// http export port  :8080
	Port string

	// interface limit
	Limit *limit.Config

	// jaeger trace config
	Trace *trace.Config
}

type RecoverInfo struct {
	Time  string      `json:"time"`
	Url   string      `json:"url"`
	Err   interface{} `json:"error"`
	Query interface{} `json:"query"`
	Stack string      `json:"stack"`
}

type CommonRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}
