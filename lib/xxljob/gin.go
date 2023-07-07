package xxljob

import (
	"github.com/gin-gonic/gin"
	"github.com/xxl-job/xxl-job-executor-go"
)

func WithGinRoute(e *gin.Engine, exec xxl.Executor) {
	if e != nil && exec != nil {
		//注册的gin的路由
		e.POST("/run", gin.WrapF(exec.RunTask))
		e.POST("/kill", gin.WrapF(exec.KillTask))
		e.POST("/log", gin.WrapF(exec.TaskLog))
		e.POST("/beat", gin.WrapF(exec.Beat))
		e.POST("/idleBeat", gin.WrapF(exec.IdleBeat))
	}
}
