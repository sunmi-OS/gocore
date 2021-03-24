package rmqv2

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/sunmi-OS/gocore/xlog"
)

type Consumer struct {
	Consumer            rocketmq.PushConsumer
	serverName          string
	messageBatchMaxSize int // default 1
	subscribeTopic      map[string]struct{}
	conf                *RocketMQConfig
	ops                 []consumer.Option
	mu                  sync.RWMutex
}

// new
func NewConsumer(conf *RocketMQConfig) (c *Consumer) {
	ops := defaultConsumerOps(conf)
	if len(conf.ConsumerOptions) > 0 {
		ops = append(ops, conf.ConsumerOptions...)
	}
	c = &Consumer{
		Consumer:       nil,
		serverName:     conf.EndPoint,
		subscribeTopic: make(map[string]struct{}),
		conf:           conf,
		ops:            ops,
	}
	return c
}

// connect to aliyun rocketmq
func (c *Consumer) Conn() (conn *Consumer, err error) {
	if c.conf.LogLevel != "" {
		rlog.SetLogLevel(string(c.conf.LogLevel))
	}
	if c.messageBatchMaxSize == 0 {
		c.messageBatchMaxSize = 1
		c.ops = append(c.ops, consumer.WithConsumeMessageBatchMaxSize(1))
	}
	newPushConsumer, err := consumer.NewPushConsumer(c.ops...)
	if err != nil {
		return nil, err
	}
	c.Consumer = newPushConsumer
	return c, nil
}

// start subscribe
func (c *Consumer) Start() (err error) {
	xlog.Infof("count [%d] start subscribe", len(c.subscribeTopic))
	return c.Consumer.Start()
}

// unsubscribe all topic
func (c *Consumer) Close() {
	if c.Consumer != nil && len(c.subscribeTopic) > 0 {
		for topic, _ := range c.subscribeTopic {
			_ = c.Consumer.Unsubscribe(topic)
			delete(c.subscribeTopic, topic)
		}
		_ = c.Consumer.Shutdown()
	}
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
		return consumer.ConsumeRetryLater, nil
	})
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.subscribeTopic[topic] = struct{}{}
	c.mu.Unlock()
	return nil
}

// 多条消息消费，需配置 client.MessageBatchMaxSize() 且size不为 1，否则不生效
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
	c.mu.Lock()
	c.subscribeTopic[topic] = struct{}{}
	c.mu.Unlock()
	return nil
}

// unsubscribe one topic
func (c *Consumer) Unsubscribe(topic string) (err error) {
	if c.Consumer == nil {
		return fmt.Errorf("[%s] is nil", c.serverName)
	}
	if err = c.Consumer.Unsubscribe(topic); err != nil {
		return err
	}
	c.mu.Lock()
	delete(c.subscribeTopic, topic)
	c.mu.Unlock()
	return
}
