package mq

import (
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

var (
	_DefaultMaxConcurrentRequests = 100
	_DefaultErrorPercentThreshold = 100
)

type Consumer struct {
	Consumer     rocketmq.PushConsumer
	ConsumerName string
	MaxCount     int // 最大次数，default = 16
	Timeout      int // 超时时间(单位：秒)，default = 60s
}

func (c *Consumer) Subscribe(topic, expression string, callback func(ctx context.Context, ext ...*primitive.MessageExt) error) (err error) {
	if c.MaxCount == 0 {
		c.MaxCount = 16
	}
	if c.Timeout == 0 {
		c.Timeout = 60
	}

	// 配置熔断
	hystrix.ConfigureCommand("rocketmq", hystrix.CommandConfig{
		Timeout:               c.Timeout * 1000,
		MaxConcurrentRequests: _DefaultMaxConcurrentRequests,
		ErrorPercentThreshold: _DefaultErrorPercentThreshold,
	})
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.ConsumerName)
	}
	err = c.Consumer.Subscribe(topic, consumer.MessageSelector{Type: consumer.TAG, Expression: expression},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			if err := callback(ctx, ext...); err != nil {
				return consumer.ConsumeRetryLater, err
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		return err
	}
	return nil
}
