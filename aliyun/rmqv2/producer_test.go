package rmqv2

import (
	"context"
	"fmt"
	"testing"
	"time"

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

	conn, err := NewProducer(conf).Conn()
	if err != nil {
		xlog.Error(err)
		return
	}

	defer conn.Close()

	for i := 0; i < 1; i++ {
		/*res,*/ err = conn.SendAsyncSingle(ctx, func(ctx context.Context, result *primitive.SendResult, err error) {
			xlog.Debugf("%+v", result)
		}, &primitive.Message{
			Topic:         "push_online",
			Body:          []byte(fmt.Sprintf("我是消息啦啦啦啦啦_%d", i)),
			TransactionId: fmt.Sprintf("TransactionId_%d", i),
		})
		if err != nil {
			xlog.Errorf("%v", err)
			return
		}
		//xlog.Debugf("res:%#v", res)
	}
	time.Sleep(time.Hour)
}
