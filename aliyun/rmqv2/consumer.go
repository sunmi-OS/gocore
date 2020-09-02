/*
	消息并发多条推送，但commit回复一次性回复，无法对单条消息进行commit回复，暂时不推荐使用
	阿里云官方推荐使用 v1.2.4 版本
*/
package rmqv2

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
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
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.ConsumerName)
	}
	err = c.Consumer.Subscribe(topic, consumer.MessageSelector{Type: consumer.TAG, Expression: expression}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
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
