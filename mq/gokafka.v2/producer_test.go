package gokafka

import (
	"context"
	"fmt"
	"testing"

	"github.com/sunmi-OS/gocore/v2/utils/closes"
)

func TestProducer(t *testing.T) {
	brokers := []string{}            // TODO: add your brokers
	topic := ""                      // TODO: add your topic
	rc := NewProducerConfig(brokers) // 注意：不能复用，每次NewProducer时都需要重新生成
	producer := NewProducer("configName", rc)
	// producer.Writer.Close()
	rc = NewProducerConfig(brokers)
	producer = NewProducer("configName", rc)
	defer closes.Close()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%v", i)
		value := []byte(fmt.Sprintf("hello kafka %v", i))
		err := producer.Send(context.Background(), topic, key, value)
		if err != nil {
			t.Fatalf("error sending message: %+v", err)
		} else {
			t.Logf("success %v\n", i)
		}
	}
}
