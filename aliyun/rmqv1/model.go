package rmqv1

import rocketmq "github.com/apache/rocketmq-client-go/core"

type RocketMQ struct {
	// 配置文件对应的 MQ client
	ConfigName string
	groupID    string
	// 设置 TCP 协议接入点，从阿里云 RocketMQ 控制台的实例详情页面获取。
	nameServer string
	// 您在阿里云账号管理控制台中创建的 AccessKeyId，用于身份认证。
	accessKey string
	// 您在阿里云账号管理控制台中创建的 AccessKeySecret，用于身份认证。
	secretKey string
	// 用户渠道，默认值为：ALIYUN。
	channel string
}

type Consumer struct {
	pushConsumer rocketmq.PushConsumer
	cc           *ConsumerConfig
}

type Producer struct {
	producer rocketmq.Producer
}

type ConsumerConfig struct {
	// 订阅 topic
	Topic string

	// expression
	Expression string

	// 最大次数，default = 16
	MaxCount int

	// 超时时间(单位：秒)，default = 60s
	Timeout int

	MessageModel rocketmq.MessageModel

	ConsumerModel rocketmq.ConsumerModel
}

type ProducerConfig struct {
	ProducerModel rocketmq.ProducerModel
}
