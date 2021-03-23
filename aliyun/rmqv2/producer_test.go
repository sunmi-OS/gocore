package rmqv2

import (
	"context"
	"fmt"
	"testing"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/sunmi-OS/gocore/xlog"
)

var ctx = context.Background()

func TestProducer(t *testing.T) {
	conf := &RocketMQConfig{
		Namespace: "MQ_INST_xxx",
		GroupName: "GID_xxx",
		EndPoint:  "http://xxx.cn-hangzhou.mq-internal.aliyuncs.com:8080",
		AccessKey: "xxx",
		SecretKey: "xxx",
		LogLevel:  LogError,
	}
	producer := NewProducer(conf)

	if err := producer.Start(); err != nil {
		xlog.Error(err)
	}

	for i := 0; i < 3; i++ {
		_, err := producer.SendSyncSingle(ctx, &primitive.Message{
			Topic:         "mdm_demo_topic",
			Body:          []byte(fmt.Sprintf("我是消息啦啦啦啦啦_%d", i)),
			TransactionId: fmt.Sprintf("TransactionId_%d", i),
		})
		if err != nil {
			xlog.Errorf("%v", err)
			return
		}
		//xlog.Debugf("%#v", result)
	}
	err := producer.Shutdown()
	if err != nil {
		xlog.Error("shutdown error:", err)
	}
}
