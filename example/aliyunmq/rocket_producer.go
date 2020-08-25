package main

import (
	"fmt"

	rocketmq "github.com/apache/rocketmq-client-go/core"
	"github.com/sunmi-OS/gocore/aliyunmq"
	"github.com/sunmi-OS/gocore/viper"
)

func main() {

	viper.C.SetDefault("aliyunmq.GroupID", "GID_xxxx")
	viper.C.SetDefault("aliyunmq.NameServer", "http://xxxxxxx.cn-hangzhou.mq-internal.aliyuncs.com:8080")
	viper.C.SetDefault("aliyunmq.AccessKey", "xxxx")
	viper.C.SetDefault("aliyunmq.SecretKey", "xxxx")
	viper.C.SetDefault("aliyunmq.Channel", "ALIYUN")

	aliyunmq.NewProducer("aliyunmq")

	// 业务中直接使用
	for i := 0; i < 1; i++ {
		msg := fmt.Sprintf("%s-%d", "Hello,Common MQ Message-", i)
		//发送消息时请设置您在阿里云 RocketMQ 控制台上申请的 Topic。
		result, err := aliyunmq.GetProducer("aliyunmq").SendMessageSync(&rocketmq.Message{Keys: "orderid", Topic: "xxxxxxx", Body: msg})
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Printf("send message: %s result: %s\n", msg, result)
	}

}
