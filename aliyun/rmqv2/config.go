/*
消息并发多条推送，但commit回复一次性回复，无法对单条消息进行commit回复，暂时不推荐使用
阿里云官方推荐使用 v1.2.4 版本
*/
package rmqv2

import (
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

const (
	LogDebug LogLevel = "debug"
	LogWarn  LogLevel = "warn"
	LogError LogLevel = "error"
	LogInfo  LogLevel = "info"
)

type LogLevel string

type RocketMQConfig struct {
	// 阿里云 实例ID
	Namespace string
	// GroupID 阿里云创建
	GroupName string
	// 设置 TCP 协议接入点，从阿里云 RocketMQ 控制台的实例详情页面获取。
	EndPoint string
	// 您在阿里云账号管理控制台中创建的 AccessKeyId，用于身份认证。
	AccessKey string
	// 您在阿里云账号管理控制台中创建的 AccessKeySecret，用于身份认证。
	SecretKey string
	// log 级别 // default info
	LogLevel LogLevel
	// 自定义消费者配置
	ConsumerOptions []consumer.Option
	// 自定义生产者配置
	ProducerOptions []producer.Option
}

func defaultConsumerOps(conf *RocketMQConfig) (ops []consumer.Option) {
	ops = []consumer.Option{
		consumer.WithNamespace(conf.Namespace),
		consumer.WithGroupName(conf.GroupName),
		consumer.WithNameServer(primitive.NamesrvAddr{conf.EndPoint}),
		consumer.WithCredentials(primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey}),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithRetry(2),
		// consumer.WithTrace(&primitive.TraceConfig{
		// 	//TraceTopic:  conf.TraceTopic, // 此处不能设置，否则消息消费commit会失效
		// 	GroupName:   conf.GroupName,
		// 	Access:      primitive.Cloud,
		// 	Resolver:    primitive.NewPassthroughResolver(primitive.NamesrvAddr{conf.EndPoint}),
		// 	Credentials: primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey},
		// }),
	}
	return ops
}

func defaultProducerOps(conf *RocketMQConfig) (ops []producer.Option) {
	ops = []producer.Option{
		producer.WithNamespace(conf.Namespace),
		producer.WithNameServer(primitive.NamesrvAddr{conf.EndPoint}),
		producer.WithCredentials(primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey}),
		producer.WithRetry(2),
	}
	// GroupName is not necessary for producer
	if conf.GroupName != "" {
		ops = append(ops, producer.WithGroupName(conf.GroupName))
	}
	return ops
}
