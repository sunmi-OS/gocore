package xxljob

import (
	"context"
	"log"
	"testing"

	"github.com/xxl-job/xxl-job-executor-go"
)

func TestNewExecutor(t *testing.T) {
	op := &Option{
		AppName:  "xxl-job-executor-sample",
		Port:     "2233",
		LogLevel: InfoLevel,
	}
	op.WithServerAddr("http://localhost:8080/xxl-job-admin")
	op.WithAccessToken("access_token_test")
	exec, err := NewExecutor(op)
	if err != nil {
		t.Error(err)
		return
	}
	exec.Init()
	exec.RegTask("demoJobHandler", func(cxt context.Context, param *xxl.RunReq) string {
		// do something
		log.Println("params: ", param)
		return "success, input what you want to write"
	})
	// change one of this two methods to start xxl-job executor

	// 1、if you want init to gin server, you can use this
	// WithGinRoute(g.Engine, exec)

	// 2、if you want init to self http server, you can use this
	// exec.Run()
}
