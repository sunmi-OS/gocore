package api

import (
	"github.com/labstack/echo"
	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/log"
	"go.uber.org/zap"
)

// 捕获panic异样防止程序终止 并且记录到日志
func ErrorLogRecover(c echo.Context) {

	if err := recover(); err != nil {
		c.Response().Write([]byte("系统错误!具体原因:" + cast.ToString(err)))
		log.Sugar.Errorw("example-log:err", zap.Error(err.(error)), map[string]interface{}{
			"URL.Path":    c.Request().URL.Path,
			"QueryParams": c.QueryParams(),
		})
	}
}
