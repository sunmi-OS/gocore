package xxljob

import (
	"errors"
	"os"

	"github.com/sunmi-OS/gocore/v2/lib/xxljob/xxl"
)

func NewExecutor(op *Option) (xxl.Executor, error) {
	// check important config
	if err := checkConfig(op); err != nil {
		return nil, err
	}
	ops := []xxl.Option{
		xxl.AccessToken(op.AccessToken), // 请求令牌
		xxl.ServerAddr(op.ServerAddr),   // xxl-job admin地址
		xxl.ExecutorPort(op.Port),       // 此处要与gin服务启动port必需一至
		xxl.RegistryKey(op.AppName),     // 执行器名称
		xxl.AppName(op.AppName),         // AppName
		xxl.SetLogLevel(op.LogLevel),    // 日志级别
	}
	if op.LogDepth > 0 {
		ops = append(ops, xxl.SetLogDepth(op.LogDepth))
	}
	executor := xxl.NewExecutor(ops...)
	return executor, nil
}

func checkConfig(op *Option) error {
	if op.AccessToken == "" {
		return errors.New("xxl-job: access token is empty")
	}
	if op.AppName == "" {
		return errors.New("xxl-job: app name is empty")
	}
	if op.Port == "" {
		op.Port = "9999"
	}
	if op.ServerAddr == "" {
		op.ServerAddr = os.Getenv("XXL_JOB_SERVER_ADDR")
		if op.ServerAddr == "" {
			return errors.New("xxl-job: server addr is empty")
		}
	}
	return nil
}
