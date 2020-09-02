/*
	消息并发多条推送，但commit回复一次性回复，无法对单条消息进行commit回复，暂时不推荐使用
	阿里云官方推荐使用 v1.2.4 版本
*/
package rmqv2

import (
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

const (
	ConsumerOPKeyGroupID           = "GroupName"
	ConsumerOPKeyNameServer        = "NameServer"
	ConsumerOPKeyCredentials       = "Credentials"
	ConsumerOPKeyMessageModel      = "MessageModel"
	ConsumerOPKeyNamespace         = "Namespace"
	ConsumerOPKeyMaxReconsumeTimes = "MaxReconsumeTimes"
	ConsumerOPKeyPullBatchSize     = "PullBatchSize"
)

type ConsumerOption struct {
	OpsMap       map[string]consumer.Option
	consumerName string
	mu           sync.RWMutex
}

func DefaultConsumerOption(consumerName string) (c *ConsumerOption) {
	conf := initConfig(consumerName)
	// default Option
	ops := map[string]consumer.Option{
		ConsumerOPKeyGroupID:           consumer.WithGroupName(conf.GroupID),
		ConsumerOPKeyNameServer:        consumer.WithNameServer(primitive.NamesrvAddr{conf.NameServer}),
		ConsumerOPKeyCredentials:       consumer.WithCredentials(primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey}),
		ConsumerOPKeyMessageModel:      consumer.WithConsumerModel(consumer.Clustering),
		ConsumerOPKeyMaxReconsumeTimes: consumer.WithMaxReconsumeTimes(16),
	}
	return &ConsumerOption{OpsMap: ops, consumerName: consumerName}
}

func (c *ConsumerOption) MessageModel(messageModel consumer.MessageModel) *ConsumerOption {
	c.mu.Lock()
	c.OpsMap[ConsumerOPKeyMessageModel] = consumer.WithConsumerModel(messageModel)
	c.mu.Unlock()
	return c
}

func (c *ConsumerOption) Namespace(namespace string) *ConsumerOption {
	c.mu.Lock()
	c.OpsMap[ConsumerOPKeyNamespace] = consumer.WithNamespace(namespace)
	c.mu.Unlock()
	return c
}

func (c *ConsumerOption) MaxReconsumeTimes(maxcount int) *ConsumerOption {
	c.mu.Lock()
	c.OpsMap[ConsumerOPKeyMaxReconsumeTimes] = consumer.WithMaxReconsumeTimes(int32(maxcount))
	c.mu.Unlock()
	return c
}

func (c *ConsumerOption) PullBatchSize(size int) *ConsumerOption {
	c.mu.Lock()
	c.OpsMap[ConsumerOPKeyPullBatchSize] = consumer.WithPullBatchSize(int32(size))
	c.mu.Unlock()
	return c
}

func (c *ConsumerOption) Start() (ps *Consumer, err error) {
	var ops []consumer.Option
	c.mu.RLock()
	for _, v := range c.OpsMap {
		ops = append(ops, v)
	}
	c.mu.RUnlock()
	pushConsumer, err := rocketmq.NewPushConsumer(ops...)
	if err != nil {
		return nil, err
	}
	if err = pushConsumer.Start(); err != nil {
		return nil, err
	}
	return &Consumer{Consumer: pushConsumer, ConsumerName: c.consumerName}, nil
}
