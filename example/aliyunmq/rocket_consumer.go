package main

import (
	"os"
	"time"
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/aliyunmq"
)

func main() {

	viper.NewConfigToToml(`
		[aliyunmq]
		NameServer = "http://xxx.cn-hangzhou.mq-internal.aliyuncs.com:8080"
 		AccessKey = "xxx"
 		SecretKey = "xxx"
 		Namespace = "xxx"
	`)

	c := aliyunmq.NewConsumer("aliyunmq", "GID_xxx", consumer.WithMaxReconsumeTimes(2))

	err := c.Subscribe("topic", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe callback: %v \n", msgs[0].MsgId)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}