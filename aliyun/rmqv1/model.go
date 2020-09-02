package rmqv1

import rocketmq "github.com/apache/rocketmq-client-go/core"

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
