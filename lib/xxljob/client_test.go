package xxljob

import (
	"context"
	"github.com/sunmi-OS/gocore/v2/lib/xxljob/xxl"
	"log"
	"testing"
)

func TestNewExecutor(t *testing.T) {
	op := &Option{
		AppName:     "xxl-job-executor-sample",
		Port:        "2233",
		AccessToken: "access_token_test",
		ServerAddr:  "http://localhost:8080/xxl-job-admin",
	}
	exec, err := NewExecutor(op)
	if err != nil {
		t.Error(err)
		return
	}
	exec.Init()
	exec.RegTask("demoJobHandler", func(cxt context.Context, param *xxl.RunReq) *xxl.ExecuteResult {
		// do something
		log.Println("params: ", param)
		return &xxl.ExecuteResult{
			Code: xxl.SuccessCode,
			Msg:  "success, input what you want to write",
		}
	})
	// change one of this two methods to start xxl-job executor

	// 1、if you want init to gin server, you can use this
	// WithGinRoute(g.Engine, exec)

	// 2、if you want init to self http server, you can use this
	// exec.Run()
}
