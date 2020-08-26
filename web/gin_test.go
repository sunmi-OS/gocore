package web

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/ecode"
)

var banner = `
                                   _
 ___   _   _   _ __    _ __ ___   (_)
/ __| | | | | | '_ \  | '_ * _ \  | |
\__ \ | |_| | | | | | | | | | | | | |
|___/  \__,_| |_| |_| |_| |_| |_| |_|
`

func TestInitGin(t *testing.T) {
	c := &Config{Port: ":2233"}
	g := InitGin(c).Release()

	initRouteG(g.Gin)

	g.Start(banner)

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

func initRouteG(g *gin.Engine) {

	g.GET("/gin/ping", func(c *gin.Context) {
		e := ecode.New(2233, "SUCCESS")
		JSON(c, nil, e)
	})
}
