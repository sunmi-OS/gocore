package rmqv2

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

type Consumer struct {
	Consumer            rocketmq.PushConsumer
	serverName          string
	messageBatchMaxSize int // default 1
	conf                *RocketMQConfig
	ops                 []consumer.Option
}

func NewConsumer(conf *RocketMQConfig) (c *Consumer) {
	ops := defaultConsumerOps(conf)
	if len(conf.ConsumerOptions) > 0 {
		ops = append(ops, conf.ConsumerOptions...)
	}
	c = &Consumer{
		Consumer:            nil,
		serverName:          conf.EndPoint,
		messageBatchMaxSize: 1,
		conf:                conf,
		ops:                 ops,
	}
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
	if c.conf.LogLevel != "" {
		rlog.SetLogLevel(string(c.conf.LogLevel))
	}
	newPushConsumer, err := consumer.NewPushConsumer(c.ops...)
	if err != nil {
		return err
	}
	c.Consumer = newPushConsumer
	return c.Consumer.Start()
}

func (c *Consumer) Unsubscribe(topic string) (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.serverName)
	}
	return c.Consumer.Unsubscribe(topic)
}

func (c *Consumer) Shutdown() (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.serverName)
	}
	return c.Consumer.Shutdown()
}

// 多条消息消费，需配置 client.ConsumeMessageBatchMaxSize() 且size不为 1，否则不生效
func (c *Consumer) SubscribeMulti(topic, expression string, callback func(ctx context.Context, ext ...*primitive.MessageExt) error) (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.serverName)
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
func (c *Consumer) SubscribeSingle(topic, expression string, callback func(ctx context.Context, ext *primitive.MessageExt) error) (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.serverName)
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
