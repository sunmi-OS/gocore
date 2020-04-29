package middleware

import (
	"fmt"

	"encoding/json"
	"runtime"
	"time"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/spf13/cast"
)

type (
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echoMiddleware.Skipper

		// Size of the stack to be printed.
		// Optional. Default value 4KB.
		StackSize int `yaml:"stack_size"`

		// DisableStackAll disables formatting stack traces of all other goroutines
		// into buffer after the trace for the current goroutine.
		// Optional. Default value false.
		DisableStackAll bool `yaml:"disable_stack_all"`

		// DisablePrintStack disables printing stack trace.
		// Optional. Default value as false.
		DisablePrintStack bool `yaml:"disable_print_stack"`
	}
	Param struct {
		Time  string      `json:"time"`
		Url   string      `json:"url"`
		Err   string      `json:"error"`
		Query interface{} `json:"query"`
		Stack string      `json:"stack"`
	}
)

var (
	// DefaultRecoverConfig is the default Recover middleware config.
	DefaultRecoverConfig = RecoverConfig{
		Skipper:           echoMiddleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
	}
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize

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
					c.Response().Write([]byte("系统错误!具体原因:" + cast.ToString(err)))
					param := &Param{
						Time:  time.Now().Format("2006-01-02 15:04:05"),
						Url:   c.Request().URL.Path,
						Err:   err.Error(),
						Query: c.QueryParams(),
						Stack: stackStr,
					}
					data, _ := json.Marshal(param)
					fmt.Printf("%s\n", string(data))

				}
			}()
			return next(c)
		}
	}
}
