package rmqv2

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type Consumer struct {
	Consumer            rocketmq.PushConsumer
	consumerName        string
	messageBatchMaxSize int // default 1
	ops                 []consumer.Option
}

func NewConsumer(conf *RocketMQConfig) (c *Consumer) {
	ops := []consumer.Option{
		consumer.WithGroupName(conf.GroupID),
		consumer.WithNameServer(primitive.NamesrvAddr{conf.NameServer}),
		consumer.WithCredentials(primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey}),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithMaxReconsumeTimes(16),
		consumer.WithConsumeMessageBatchMaxSize(1),
	}

	if len(conf.Options) > 0 {
		ops = append(ops, conf.Options...)
	}
	c = &Consumer{
		Consumer:            nil,
		consumerName:        conf.NameServer,
		messageBatchMaxSize: 1,
		ops:                 ops,
	}
	return c
}

func (c *Consumer) Namespace(namespace string) *Consumer {
	c.ops = append(c.ops, consumer.WithNamespace(namespace))
	return c
}

func (c *Consumer) MessageModel(messageModel consumer.MessageModel) *Consumer {
	c.ops = append(c.ops, consumer.WithConsumerModel(messageModel))
	return c
}

func (c *Consumer) PullBatchSize(size int) *Consumer {
	c.ops = append(c.ops, consumer.WithPullBatchSize(int32(size)))
	return c
}

func (c *Consumer) ConsumeMessageBatchMaxSize(size int) *Consumer {
	c.messageBatchMaxSize = size
	c.ops = append(c.ops, consumer.WithConsumeMessageBatchMaxSize(size))
	return c
}

func (c *Consumer) Start() (err error) {
	newPushConsumer, err := consumer.NewPushConsumer(c.ops...)
	if err != nil {
		return err
	}
	c.Consumer = newPushConsumer
	return nil
}

// 多条消息消费，需配置 client.ConsumeMessageBatchMaxSize() 且size不为 1，否则不生效
func (c *Consumer) SubscribeMultiMessage(topic, expression string, callback func(ctx context.Context, ext ...*primitive.MessageExt) error) (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.consumerName)
	}
	selector := consumer.MessageSelector{Type: consumer.TAG, Expression: expression}
	err = c.Consumer.Subscribe(topic, selector, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		// 多条消息
		if err = callback(ctx, ext...); err != nil {
			return consumer.ConsumeRetryLater, err
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 单条消息消费 default
func (c *Consumer) SubscribeSingleMessage(topic, expression string, callback func(ctx context.Context, ext *primitive.MessageExt) error) (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.consumerName)
	}
	selector := consumer.MessageSelector{Type: consumer.TAG, Expression: expression}
	err = c.Consumer.Subscribe(topic, selector, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		// 单条消息
		if c.messageBatchMaxSize == 1 {
			if err = callback(ctx, ext[0]); err != nil {
				return consumer.ConsumeRetryLater, err
			}
			return consumer.ConsumeSuccess, nil
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Consumer) Unsubscribe(topic string) (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.consumerName)
	}
	return c.Consumer.Unsubscribe(topic)
}

func (c *Consumer) Shutdown() (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.consumerName)
	}
	return c.Consumer.Shutdown()
}
