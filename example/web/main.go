package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/ecode"
	"github.com/sunmi-OS/gocore/web"
)

func main() {
	c := &web.Config{Port: ":2233"}
	g := web.InitGin(c)
	g.Gin.Use(g.CORS())
	g.Gin.Use(g.Recovery())
	g.Release() // release or not release

	initRouteG(g.Gin)
	g.Start()

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
		e := ecode.New(2233, "success")
		web.JSON(c, nil, e)
	})

	g.GET("/gin/file", func(c *gin.Context) {
		//err := ecode.New(2323, "asdsda")
		web.File(c, "echo_test.go", "gin.go")
	})

	g.GET("/gin/page", func(c *gin.Context) {
		e := ecode.New(2233, "success")

		rsp := []string{"1", "2", "3", "4", "5"}
		web.JSON(c, web.Pager{PageNo: 1, PageSize: 5}.Apply(10, rsp), e)
	})
}
