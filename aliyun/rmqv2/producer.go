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

func (p *Producer) Start() (err error) {
	if p.conf.LogLevel != "" {
		rlog.SetLogLevel(string(p.conf.LogLevel))
	}
	defaultProducer, err := producer.NewDefaultProducer(p.ops...)
	if err != nil {
		return err
	}
	p.Producer = defaultProducer
	return p.Producer.Start()
}

func (p *Producer) Shutdown() (err error) {
	if p.Producer == nil {
		return fmt.Errorf("[%s] is nil", p.serverName)
	}

	return p.Producer.Shutdown()
}

func (p *Producer) SendSyncSingle(c context.Context, message *primitive.Message) (result *primitive.SendResult, err error) {
	if p.Producer == nil {
		return nil, fmt.Errorf("[%s] is nil", p.serverName)
	}
	return p.Producer.SendSync(c, message)
}

func (p *Producer) SendSyncMulti(c context.Context, messages []*primitive.Message) (result *primitive.SendResult, err error) {
	if p.Producer == nil {
		return nil, fmt.Errorf("[%s] is nil", p.serverName)
	}
	return p.Producer.SendSync(c, messages...)
}

// 不能用，暂时注释掉
//func (p *Producer) SendAsyncSingle(c context.Context, message *primitive.Message) (result *primitive.SendResult, err error) {
//	if p.Producer == nil {
//		return nil, fmt.Errorf("[%s] is nil", p.serverName)
//	}
//
//	err = p.Producer.SendAsync(c, func(ctx context.Context, res *primitive.SendResult, err error) {
//		if err != nil {
//			//result = res
//			xlog.Debugf("%#v", res)
//		}
//	}, message)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}
//
//func (p *Producer) SendAsyncMulti(c context.Context, messages []*primitive.Message) (result *primitive.SendResult, err error) {
//	if p.Producer == nil {
//		return nil, fmt.Errorf("[%s] is nil", p.serverName)
//	}
//	err = p.Producer.SendAsync(c, func(ctx context.Context, res *primitive.SendResult, err error) {
//		if err != nil {
//			result = res
//		}
//	}, messages...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}

// 不能用，暂时注释掉
//func (p *Producer) SendOneWaySingle(c context.Context, message *primitive.Message) (err error) {
//	if p.Producer == nil {
//		return fmt.Errorf("[%s] is nil", p.serverName)
//	}
//	return p.Producer.SendOneWay(c, message)
//}
//
//func (p *Producer) SendOneWayMulti(c context.Context, messages []*primitive.Message) (err error) {
//	if p.Producer == nil {
//		return fmt.Errorf("[%s] is nil", p.serverName)
//	}
//	return p.Producer.SendOneWay(c, messages...)
//}
