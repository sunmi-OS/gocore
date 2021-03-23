package rmqv2

import (
	"context"
	"testing"

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
	}
	consumer := NewConsumer(conf)

	if err := consumer.Start(); err != nil {
		xlog.Error(err)
	}
	defer consumer.Shutdown()

	err := consumer.SubscribeSingle("mdm_demo_topic", "*", func(c context.Context, ext *primitive.MessageExt) error {
		xlog.Debugf("Message:%#v", ext)
		xlog.Debugf("body:%v", string(ext.Body))
		return nil
	})
	if err != nil {
		xlog.Error(err)
		return
	}
	for {

	}
}
