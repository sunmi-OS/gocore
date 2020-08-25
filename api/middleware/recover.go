package middleware

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
)

type Param struct {
	Time  string      `json:"time"`
	Url   string      `json:"url"`
	Err   string      `json:"error"`
	Query interface{} `json:"query"`
	Stack string      `json:"stack"`
}

var returnMsg = `{"message":"Internal Server Error"}`

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover(msg ...string) echo.MiddlewareFunc {
	if len(msg) > 0 {
		returnMsg = msg[0]
	}
	return RecoverWithConfig(echoMiddleware.DefaultRecoverConfig)
}

func SetReturnMsg(msg string) {
	returnMsg = msg
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config echoMiddleware.RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = echoMiddleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = echoMiddleware.DefaultRecoverConfig.StackSize

	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					stackStr := ""
					length := runtime.Stack(stack, false)
					if !config.DisablePrintStack {
						stackStr = string(stack[:length])
					}
					c.Response().Header().Set("Content-Type", "application/json;charset=UTF-8")
					c.Response().Write([]byte(returnMsg))
					param := &Param{
						Time:  time.Now().Format("2006-01-02 15:04:05"),
						Url:   c.Request().URL.Path,
						Err:   err.Error(),
						Query: c.QueryParams(),
						Stack: stackStr,
					}
					data, _ := json.Marshal(param)
					fmt.Println(string(data))

				}
			}()
			return next(c)
		}
	}
}
