package xxljob

import (
	"errors"
	"os"

	"github.com/xxl-job/xxl-job-executor-go"
)

func NewExecutor(op *Option) (xxl.Executor, error) {
	// check important config
	if err := checkConfig(op); err != nil {
		return nil, err
	}
	executor := xxl.NewExecutor(
		xxl.AccessToken(op.accessToken), // 请求令牌
		xxl.ServerAddr(op.serverAddr),   // xxl-job admin地址
		xxl.ExecutorPort(op.Port),       // 此处要与gin服务启动port必需一至
		xxl.RegistryKey(op.AppName),     // 执行器名称
		xxl.SetLogger(newLogger(op.LogLevel)),
	)
	return executor, nil
}

func checkConfig(op *Option) error {
	if op.LogLevel == "" {
		op.LogLevel = InfoLevel
	}
	if op.accessToken == "" {
		return errors.New("xxl-job: access token is empty")
	}
	if op.AppName == "" {
		return errors.New("xxl-job: app name is empty")
	}
	if op.Port == "" {
		op.Port = "9999"
	}
	if op.serverAddr == "" {
		op.serverAddr = os.Getenv("XXL_JOB_SERVER_ADDR")
		if op.serverAddr == "" {
			return errors.New("xxl-job: server addr is empty")
		}
	}
	return nil
}
