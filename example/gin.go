package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/api"
	"github.com/sunmi-OS/gocore/v2/utils/closes"
)

func main() {
	closes.AddShutdown(closes.ModuleClose{
		Name:     "aaa",
		Priority: 0,
		Func: func() {
			// do something
			log.Println("aaa close")
		},
	})

	hs := api.NewGinServer(
		//api.WithServerHost(""),
		api.WithServerPort(2233),
		api.WithServerDebug(true),
		api.WithServerTimeout(time.Second*10),
		api.WithOpenTrace(false),
	)
	// init route
	initRoute(hs.Gin)

	// add close hook
	hs.AddShutdownHook(func(c context.Context) {
		// do something when http server showdown
	})
	// start server
	hs.Start()
}

func initRoute(g *gin.Engine) {
	group := g.Group("/api")
	{
		group.POST("/test", func(c *gin.Context) {
			newContext := api.NewContext(c)

			// do something

			// return
			newContext.RetJSON("response data", nil)
		})
	}
}
