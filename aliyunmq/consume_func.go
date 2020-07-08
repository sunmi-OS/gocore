package aliyunmq

import (
	"fmt"
	rocketmq "github.com/apache/rocketmq-client-go/core"
)

type ConsumeFunc struct {
	MaxCount int // 最大16次
}

func NewConsumeFunc() ConsumeFunc {

	return ConsumeFunc{
		MaxCount: 16,
	}

}

// 最大重试次数，超过次数发邮件报警等功能可以直接扩展
func (c ConsumeFunc) Middleware(f func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus) func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {

	return func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {


		if msg.ReconsumeTimes >= c.MaxCount {
			return rocketmq.ConsumeSuccess
		}

		ch := make(chan rocketmq.ConsumeStatus)

		defer close(ch)

		go func(f func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus, ch chan rocketmq.ConsumeStatus) {

			defer func() {
				if err := recover(); err != nil {
					ch <- rocketmq.ReConsumeLater
					fmt.Println(err)
				}
			}()
			ch <- f(msg)
		}(f, ch)

		status := <- ch
		return status
	}
}
