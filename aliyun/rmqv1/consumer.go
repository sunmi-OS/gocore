package rmqv1

import (
	"errors"
	"log"

	"github.com/afex/hystrix-go/hystrix"
	rocketmq "github.com/apache/rocketmq-client-go/core"
)

var (
	_DefaultMaxConcurrentRequests = 100
	_DefaultErrorPercentThreshold = 100
	_hystrixName                  = "rmqv1"
)

func GetPushConsumer(configName string) (consumer rocketmq.PushConsumer, ok bool) {
	pc, ok := ConsumerList.Load(configName)
	if ok {
		return pc.(rocketmq.PushConsumer), ok
	}
	return nil, false
}

func (c *Consumer) Subscribe(callback func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus) (err error) {
	if c.pushConsumer == nil {
		return errors.New("consumer is nil")
	}
	var (
		timeout    = 60
		expression = "*"
		maxCount   = 16
	)
	if c.cc.Expression != "" {
		expression = c.cc.Expression
	}
	if c.cc.MaxCount != 0 {
		maxCount = c.cc.MaxCount
	}
	if c.cc.Timeout != 0 {
		timeout = c.cc.Timeout
	}

	hystrix.ConfigureCommand(_hystrixName, hystrix.CommandConfig{
		Timeout:               timeout * 1000,
		MaxConcurrentRequests: _DefaultMaxConcurrentRequests,
		ErrorPercentThreshold: _DefaultErrorPercentThreshold,
	})

	return c.pushConsumer.Subscribe(c.cc.Topic, expression, func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {
		if msg.ReconsumeTimes >= maxCount {
			return rocketmq.ConsumeSuccess
		}
		ch := make(chan rocketmq.ConsumeStatus)
		e := hystrix.Go(_hystrixName, func() error {
			defer func() {
				if err := recover(); err != nil {
					ch <- rocketmq.ReConsumeLater
				}
				close(ch)
			}()
			ch <- callback(msg)
			return nil
		}, nil)

		select {
		case out := <-ch:
			// success
			return out
		case err := <-e:
			// failure
			log.Println("hystrix Error:", err)
			return rocketmq.ReConsumeLater
		}
	})
}

func (c *Consumer) Start() error {
	if c.pushConsumer != nil {
		return c.pushConsumer.Start()
	}
	return errors.New("consumer is nil")
}

func (c *Consumer) Close() error {
	if c.pushConsumer != nil {
		return c.pushConsumer.Shutdown()
	}
	return errors.New("consumer is nil")
}
