package ecode

import "sync"

var (
	ErrorMap = new(sync.Map)
	// base error
	OK         = New(0, "SUCCESS")
	RequestErr = New(10400, "请求参数错误，请检查。")
	ServerErr  = New(11503, "服务异常，请稍后再试。")
)
