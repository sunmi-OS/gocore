package resutful

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

// 捕获panic异样防止程序终止 并且记录到日志
func ErrorLogRecover(c echo.Context) {

	if err := recover(); err != nil {
		_, err2 := c.Response().Write([]byte("系统错误!具体原因:" + cast.ToString(err)))
		log.Println("example-log:err", err2, map[string]interface{}{
			"URL.Path":    c.Request().URL.Path,
			"QueryParams": c.QueryParams(),
		})

	}
}
