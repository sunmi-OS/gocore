package rmqv2

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

type Producer struct {
	Producer   rocketmq.Producer
	serverName string
	conf       *RocketMQConfig
	ops        []producer.Option
}

func NewProducer(conf *RocketMQConfig) (p *Producer) {
	ops := defaultProducerOps(conf)
	if len(conf.ProducerOptions) > 0 {
		ops = append(ops, conf.ProducerOptions...)
	}
	p = &Producer{
		Producer:   nil,
		serverName: conf.EndPoint,
		conf:       conf,
		ops:        ops,
	}
	return p
}

// connect to aliyun rocketmq
func (p *Producer) Conn() (conn *Producer, err error) {
	if p.conf.LogLevel != "" {
		rlog.SetLogLevel(string(p.conf.LogLevel))
	}
	defaultProducer, err := producer.NewDefaultProducer(p.ops...)
	if err != nil {
		return nil, err
	}
	p.Producer = defaultProducer
	if err = p.Producer.Start(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Producer) Close() {
	if p.Producer != nil {
		_ = p.Producer.Shutdown()
	}
}

// 同步单条消息发送，对应消费 topic 的 MessageBatchMaxSize = 1时用
func (p *Producer) SendSyncSingle(c context.Context, message *primitive.Message) (result *primitive.SendResult, err error) {
	if p.Producer == nil {
		return nil, fmt.Errorf("[%s] is nil", p.serverName)
	}
	return p.Producer.SendSync(c, message)
}

// 异步单条消息发送，对应消费 topic 的 MessageBatchMaxSize = 1时用
func (p *Producer) SendAsyncSingle(c context.Context, callback func(ctx context.Context, result *primitive.SendResult, err error), message *primitive.Message) (err error) {
	if p.Producer == nil {
		return fmt.Errorf("[%s] is nil", p.serverName)
	}
	if callback == nil {
		callback = func(ctx context.Context, result *primitive.SendResult, err error) {}
	}
	err = p.Producer.SendAsync(c, callback, message)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) SendOneWaySingle(c context.Context, message *primitive.Message) (err error) {
	if p.Producer == nil {
		return fmt.Errorf("[%s] is nil", p.serverName)
	}
	return p.Producer.SendOneWay(c, message)
}
