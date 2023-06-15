package gokafka

import (
	"context"
	"testing"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils/closes"

	"github.com/segmentio/kafka-go"
)

func TestConsumer(t *testing.T) {
	brokers := []string{}                            // TODO: add your brokers
	groupID := ""                                    // TODO: add your groupID
	topic := ""                                      // TODO: add your topic
	rc := NewConsumerConfig(brokers, groupID, topic) // 注意：rc不要复用，每次NewConsumer时都需要重新生成
	consumer := NewConsumer(rc)

	go func() {
		glog.InfoF("start consumer, %#v", rc)
		err := consumer.Handle(context.Background(), func(msg kafka.Message) error {
			t.Logf("key=%s, value=%s\n", msg.Key, msg.Value)
			return nil
		})
		if err != nil && err != context.Canceled {
			t.Errorf("err=%+v", err)
		}
	}()
	time.Sleep(time.Second * 10)
	defer closes.Close()
}
