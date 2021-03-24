package rmqv2

import (
	"context"
	"testing"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/sunmi-OS/gocore/xlog"
)

func TestConsumer(t *testing.T) {
	conf := &RocketMQConfig{
		Namespace: "MQ_INST_xxx",
		GroupName: "GID_xxx",
		EndPoint:  "http://xxx.cn-hangzhou.mq-internal.aliyuncs.com:8080",
		AccessKey: "xxx",
		SecretKey: "xxx",
		LogLevel:  LogError,
		// 自定义配置
		//ConsumerOptions: []consumer.Option{
		//	consumer.WithMaxReconsumeTimes(10),
		//	consumer.WithConsumeMessageBatchMaxSize(1),
		//	consumer.WithPullInterval(time.Millisecond),
		//	consumer.WithPullBatchSize(10),
		//},
	}
	conn, err := NewConsumer(conf).Conn()
	if err != nil {
		xlog.Error(err)
		return
	}
	defer conn.Close()

	if err = conn.SubscribeSingle("mdm_demo_topic", "*", func(c context.Context, ext *primitive.MessageExt) error {
		xlog.Debugf("body:%v", string(ext.Body))
		return nil
	}); err != nil {
		xlog.Error(err)
		return
	}

	if err = conn.Start(); err != nil {
		xlog.Error(err)
	}

	time.Sleep(time.Hour)
}
