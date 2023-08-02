package example

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/api"
)

func ExampleNewGinServer() {
	hs := api.NewGinServer(
		//api.WithServerHost(""),
		api.WithServerPort(2233),
		api.WithServerDebug(true),
		api.WithServerTimeout(time.Second*30),
		api.WithOpenTrace(false),
	)
	// init route
	initRoute(hs.Gin)

	// add close hook
	hs.AddCloseHook(func(c context.Context) {
		// do something when server close
	})
	// add exit hook
	hs.AddExitHook(func(c context.Context) {
		// do something when process exit
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
