package aliyunmq

import (
	"github.com/afex/hystrix-go/hystrix"
	rocketmq "github.com/apache/rocketmq-client-go/core"
	"log"
)

type ConsumeFunc struct {
	MaxCount int // 最大16次
	TimeOut  int // 超时时间
}

func NewConsumeFunc() ConsumeFunc {

	return ConsumeFunc{
		MaxCount: 16,
		TimeOut:  60,
	}
}

func (c *ConsumeFunc) SetMaxCount(i int) {
	c.MaxCount = i
}

func (c *ConsumeFunc) SetTimeOut(i int) {
	c.TimeOut = i
}

// 最大重试次数，超过次数发邮件报警等功能可以直接扩展
func (c *ConsumeFunc) Middleware(f func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus) func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {

	hystrix.ConfigureCommand("rocketmq", hystrix.CommandConfig{
		Timeout:               c.TimeOut * 1000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 100,
	})

	return func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {

		if msg.ReconsumeTimes >= c.MaxCount {
			return rocketmq.ConsumeSuccess
		}
		ch := make(chan rocketmq.ConsumeStatus)
		errors := hystrix.Go("rocketmq", func() error {
			defer func() {
				if err := recover(); err != nil {
					ch <- rocketmq.ReConsumeLater
				}
				close(ch)
			}()
			ch <- f(msg)
			return nil
		}, nil)

		select {
		case out := <-ch:
			// success
			return out
		case err := <-errors:
			// failure
			log.Println("hystrix Error:", err)
			return rocketmq.ReConsumeLater
		}
	}
}
