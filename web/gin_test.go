package web

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/ecode"
)

var banner = `
 __ _                      _           
/ _* |   ___     ___    __| |    ___
\__, |  / _ \   / _ \  / _* |   (_-<   
|___/   \___/   \___/  \__,_|   /__/_ 

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

			// todo something

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
