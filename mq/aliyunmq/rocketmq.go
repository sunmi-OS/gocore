package aliyunmq

import (
	"errors"
	"runtime"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
)

// RocketMQConfig 基础配置
type RocketMQConfig struct {
	// 设置 TCP 协议接入点，从阿里云 RocketMQ 控制台的实例详情页面获取。
	NameServer string
	// 命名空间（阿里云上的实例ID）
	Namespace string
	// 您在阿里云账号管理控制台中创建的 AccessKeyId，用于身份认证。
	AccessKey string
	// 您在阿里云账号管理控制台中创建的 AccessKeySecret，用于身份认证。
	SecretKey string
	// 是否自建RocketMQ true-自建 false-阿里云托管版。
	IsLocal bool
}

// initConfig Initial configuration
func initConfig(configName string) RocketMQConfig {
	mqConfig := RocketMQConfig{
		NameServer: viper.GetEnvConfig(configName + ".NameServer").String(),
		AccessKey:  viper.GetEnvConfig(configName + ".AccessKey").String(),
		SecretKey:  viper.GetEnvConfig(configName + ".SecretKey").String(),
		Namespace:  viper.GetEnvConfig(configName + ".Namespace").String(),
		IsLocal:    viper.GetEnvConfig(configName + ".IsLocal").Bool(),
	}
	err := checkConfig(mqConfig)
	if err != nil {
		panic(err)
	}
	// Set the default log level to error
	LogError()
	// Set the panic handler
	primitive.PanicHandler = func(i interface{}) {
		stacktrace := stack()
		glog.ErrorF("rocketmq panic recovered:%+v, stacktrace:%s", i, stacktrace)
	}
	return mqConfig
}

func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		bufSize := len(buf)
		if n < bufSize {
			return buf[:n]
		}
		buf = make([]byte, bufSize*2)
	}
}

// checkConfig Check configuration integrity
func checkConfig(conf RocketMQConfig) (err error) {
	if conf.AccessKey == "" || conf.NameServer == "" || conf.SecretKey == "" {
		err = errors.New("missing required configuration")
		return
	}
	if conf.IsLocal && conf.Namespace != "" {
		err = errors.New("namespace must be empty when islocal is true")
		return
	}
	if !conf.IsLocal && conf.Namespace == "" {
		err = errors.New("namespace can not be empty")
		return
	}
	return
}
