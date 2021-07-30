package api

import (
	"time"

	"github.com/fvbock/endless"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		time.Sleep(10 * time.Second)

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	endless.ListenAndServe(":8080", r) // 监听并在 0.0.0.0:8080 上启动服务

}
