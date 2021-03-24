package main


import (
	"context"
	"time"
	"fmt"

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

	aliyunmq.NewProducer("aliyunmq")

	p := aliyunmq.GetProducer("aliyunmq")

	for i := 0; i < 10; i++ {
		err := p.SendAsync(context.Background(),
			func(ctx context.Context, result *primitive.SendResult, e error) {
				if e != nil {
					fmt.Printf("receive message error: %s\n", e)
				} else {
					fmt.Printf("send message success: result=%s\n", result.String())
				}
			}, primitive.NewMessage("topic", []byte("Hello RocketMQ Go Client!")))

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		}
	}

	time.Sleep(1 * time.Hour)
	err := p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}
