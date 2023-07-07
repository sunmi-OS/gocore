package xxljob

import "github.com/sunmi-OS/gocore/v2/lib/xxljob/xxl"

type Option struct {
	AppName     string       // xxl-job executor's app name
	Port        string       // the port of xxl-job executor, default is 9999
	AccessToken string       // access token of xxl-job executor
	ServerAddr  string       // the address of xxl-job admin, if empty, will use env XXL_JOB_SERVER_ADDR
	LogLevel    xxl.LogLevel // 日志级别
	LogDepth    int          // 日志深度
}

// ExecuteSuccess 执行成功返回
func ExecuteSuccess(msg string) *xxl.ExecuteResult {
	return &xxl.ExecuteResult{
		Code: xxl.SuccessCode,
		Msg:  msg,
	}
}

// ExecuteFail 执行失败返回
func ExecuteFail(msg string) *xxl.ExecuteResult {
	return &xxl.ExecuteResult{
		Code: xxl.FailureCode,
		Msg:  msg,
	}
}
