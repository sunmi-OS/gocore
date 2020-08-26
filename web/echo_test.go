package web

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestInitEcho(t *testing.T) {
	c := &Config{Port: ":2234"}
	e := InitEcho(c).Release()

	initRouteE(e.Echo)

	e.Start(banner)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			// todo something

			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func initRouteE(e *echo.Echo) {
	e.GET("/echo/ping", func(c echo.Context) error {
		//err := ecode.New(2323, "asdsda")
		return c.JSON(http.StatusOK, CommonRsp{
			Code:    2234,
			Message: "SUCCESS",
			Data:    nil,
		})
	})
}
