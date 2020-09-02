package rmqv1

import (
	"errors"

	rocketmq "github.com/apache/rocketmq-client-go/core"
)

func GetProducer(configName string) (producer rocketmq.Producer, ok bool) {
	p, ok := ProducerList.Load(configName)
	if ok {
		return p.(rocketmq.Producer), ok
	}
	return nil, false
}

func (p *Producer) SendMessageSync(msg *rocketmq.Message) (*rocketmq.SendResult, error) {
	if p.producer != nil {
		return p.producer.SendMessageSync(msg)
	}
	return nil, errors.New("producer is nil")
}

func (p *Producer) SendMessageOneway(msg *rocketmq.Message) error {
	if p.producer != nil {
		return p.producer.SendMessageOneway(msg)
	}
	return errors.New("producer is nil")
}

func (p *Producer) SendMessageOrderly(msg *rocketmq.Message, selector rocketmq.MessageQueueSelector, arg interface{}, autoRetryTimes int) (*rocketmq.SendResult, error) {
	if p.producer != nil {
		return p.producer.SendMessageOrderly(msg, selector, arg, autoRetryTimes)
	}
	return nil, errors.New("producer is nil")
}

func (p *Producer) SendMessageOrderlyByShardingKey(msg *rocketmq.Message, shardingkey string) (*rocketmq.SendResult, error) {
	if p.producer != nil {
		return p.producer.SendMessageOrderlyByShardingKey(msg, shardingkey)
	}
	return nil, errors.New("producer is nil")

}
func (p *Producer) Start() error {
	if p.producer != nil {
		return p.producer.Start()
	}
	return errors.New("producer is nil")
}

func (p *Producer) Close() error {
	if p.producer != nil {
		return p.producer.Shutdown()
	}
	return errors.New("producer is nil")
}
