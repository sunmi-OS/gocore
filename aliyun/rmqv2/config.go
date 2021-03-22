/*
	消息并发多条推送，但commit回复一次性回复，无法对单条消息进行commit回复，暂时不推荐使用
	阿里云官方推荐使用 v1.2.4 版本
*/
package rmqv2

import (
	"errors"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/sunmi-OS/gocore/viper"
)

type RocketMQConfig struct {
	GroupID string
	// 设置 TCP 协议接入点，从阿里云 RocketMQ 控制台的实例详情页面获取。
	NameServer string
	// 您在阿里云账号管理控制台中创建的 AccessKeyId，用于身份认证。
	AccessKey string
	// 您在阿里云账号管理控制台中创建的 AccessKeySecret，用于身份认证。
	SecretKey string
	// 用户渠道，默认值为：ALIYUN。
	Channel string
	// 配置
	Options []consumer.Option
}

func initConfig(configName string) (config *RocketMQConfig) {

	config = &RocketMQConfig{
		GroupID:    viper.GetEnvConfig(configName + ".GroupID"),
		NameServer: viper.GetEnvConfig(configName + ".NameServer"),
		AccessKey:  viper.GetEnvConfig(configName + ".AccessKey"),
		SecretKey:  viper.GetEnvConfig(configName + ".SecretKey"),
		Channel:    viper.GetEnvConfig(configName + ".Channel"),
	}

	err := checkConfig(config)
	if err != nil {
		panic(err)
	}
	return
}

func checkConfig(conf *RocketMQConfig) (err error) {
	if conf == nil {
		return errors.New("nil RocketMQConfig")
	}
	if conf.AccessKey == "" || conf.Channel == "" || conf.GroupID == "" || conf.NameServer == "" || conf.SecretKey == "" {
		err = errors.New("config Missing parameter")
	}
	return
}
