package aliyunmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// NewConsumer 初始化消费者
func NewConsumer(configName, groupID string, option ...consumer.Option) rocketmq.PushConsumer {

	conf := initConfig(configName)

	consumerOption := []consumer.Option{
		consumer.WithGroupName(groupID),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.NameServer})),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: conf.AccessKey,
			SecretKey: conf.SecretKey,
		}),
		consumer.WithConsumeMessageBatchMaxSize(1), // 批量获取数量 默认每次获取一条处理一条
		consumer.WithNamespace(conf.Namespace),
		consumer.WithTrace(&primitive.TraceConfig{
			GroupName:   groupID,
			Access:      primitive.Cloud,
			Resolver:    primitive.NewPassthroughResolver(primitive.NamesrvAddr{conf.NameServer}),
			Credentials: primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey},
		}),
	}

	// 支持自定义配置
	if len(option) > 0 {
		consumerOption = append(consumerOption, option...)
	}

	conn, err := rocketmq.NewPushConsumer(consumerOption...)

	if err != nil {
		panic(err)
	}

	err = conn.Start()
	if err != nil {
		panic(err)
	}

	return conn
}
