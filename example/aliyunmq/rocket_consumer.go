package main

import (
	"github.com/sunmi-OS/gocore/aliyunmq"
	"fmt"
	rocketmq "github.com/apache/rocketmq-client-go/core"
	"github.com/sunmi-OS/gocore/viper"
	"time"
)

func main() {

	viper.C.SetDefault("aliyunmq.GroupID", "GID_xxxx")
	viper.C.SetDefault("aliyunmq.NameServer", "http://xxxx.cn-hangzhou.mq-internal.aliyuncs.com:8080")
	viper.C.SetDefault("aliyunmq.AccessKey", "xxxx")
	viper.C.SetDefault("aliyunmq.SecretKey", "xxxx")
	viper.C.SetDefault("aliyunmq.Channel", "ALIYUN")

	aliyunmq.NewConsumer("aliyunmq")

	consumer, err := aliyunmq.GetPushConsumer("aliyunmq")
	if err != nil {
		panic(err)
	}

	consumeFunc := aliyunmq.NewConsumeFunc()

	defer func() {
		err = consumer.Shutdown()
		if err != nil {
			println("consumer shutdown failed")
			return
		}
	}()

	consumer.Subscribe("xxxxx", "*", consumeFunc.Middleware(func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {


		fmt.Printf("A message received, MessageID:%s, Body:%s \n", msg.MessageID, msg.Body)
		fmt.Println(msg.ReconsumeTimes)

		//消费成功回复 ConsumeSuccess，消费失败回复 ReConsumeLater。此时会触发消费重试。
		return rocketmq.ConsumeSuccess
	}))

	err = consumer.Start()
	if err != nil {
		println("consumer start failed,", err)
		return
	}

	fmt.Printf("consumer: %s started...\n", consumer)

	ch := make(chan interface{})
	<-ch
}
