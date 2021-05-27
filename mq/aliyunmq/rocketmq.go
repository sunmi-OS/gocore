package aliyunmq

import (
	"errors"
	viper2 "github.com/sunmi-OS/gocore/v2/conf/viper"
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
}

// initConfig 通过viper初始化配置
func initConfig(configName string) RocketMQConfig {

	mqConfig := RocketMQConfig{
		NameServer: viper2.GetEnvConfig(configName + ".NameServer"),
		AccessKey:  viper2.GetEnvConfig(configName + ".AccessKey"),
		SecretKey:  viper2.GetEnvConfig(configName + ".SecretKey"),
		Namespace:  viper2.GetEnvConfig(configName + ".Namespace"),
	}
	err := checkConfig(mqConfig)
	if err != nil {
		panic(err)
	}
	// 默认日志等级 Error
	LogError()

	return mqConfig
}

// checkConfig 检查配置完整性
func checkConfig(conf RocketMQConfig) (err error) {

	if conf.AccessKey == "" || conf.Namespace == "" || conf.NameServer == "" || conf.SecretKey == "" {
		err = errors.New("config Missing parameter")
	}
	return
}
