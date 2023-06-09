package gokafka

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils/closes"

	"github.com/segmentio/kafka-go"
)

func TestPrefix(t *testing.T) {
	envstr := os.Getenv("RUN_TIME")
	os.Setenv("RUN_TIME", "dev")

	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello", "Dev-Hello"},
		{"", "Dev-"},
		{"this is a long string", "Dev-this is a long string"},
	}

	for _, tc := range testCases {
		result := RunTimePrefix(tc.input)
		if result != tc.expected {
			t.Errorf("Unexpected result for input %v: got %v, want %v", tc.input, result, tc.expected)
		}

		result = PreTopicPrefix(tc.input)
		if result != tc.expected {
			t.Errorf("Unexpected result for input %v: got %v, want %v", tc.input, result, tc.expected)
		}
	}
	err := os.Setenv("RUN_TIME", "pre")
	if err != nil {
		t.Errorf("set RUN_TIME failed err=%v", err)
	}
	result := PreTopicPrefix("hello")
	if result != "Onl-hello" {
		t.Errorf("Unexpected result for input %v: got %v, want %v", "hello", result, "Onl-hello")
	}

	os.Setenv("RUN_TIME", envstr)
}

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
