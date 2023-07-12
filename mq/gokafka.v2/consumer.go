package gokafka

import (
	"context"
	"io"
	"time"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils/closes"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader *kafka.Reader
	ctx    context.Context
	cancel context.CancelFunc
}

func NewConsumerConfig(brokers []string, groupID string, topic string) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, //10K
		MaxBytes:       10e6, //10MB
		MaxWait:        time.Second,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	}
}

func NewVipConsumerConfig(brokername string, groupID string, topic string) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:        viper.GetEnvConfig(brokername + ".Brokers").SliceString(),
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3, //10K
		MaxBytes:       10e6, //10MB
		MaxWait:        time.Second,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	}
}

// NewConsumer conf每次重新生成
func NewConsumer(conf kafka.ReaderConfig) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Consumer{
		ctx:    ctx,
		cancel: cancel,
		Reader: kafka.NewReader(conf),
	}
	closes.AddShutdown(closes.ModuleClose{
		Name:     "Kafka Consumer Close",
		Priority: closes.MQPriority,
		Func: func() {
			_ = c.Close()
		},
	})
	return c
}

func (kr *Consumer) Handle(ctx context.Context, handle func(msg kafka.Message) error) error {
	for {
		select {
		case <-ctx.Done():
			glog.InfoF("Kafka Consumer.Handle ctx done")
			return ctx.Err()
		case <-kr.ctx.Done():
			glog.InfoF("Kafka Consumer.Handle kr.ctx done")
			return kr.ctx.Err()
		default:
			m, err := kr.Reader.FetchMessage(ctx)
			// io.EOF means consumer closed
			// io.ErrClosedPipe means committing messages on the consumer,
			// kafka will refire the messages on uncommitted messages, ignore
			if err == io.EOF || err == io.ErrClosedPipe {
				glog.InfoF("Kafka Consumer.FetchMessage error:%v(the reader has been closed)", err)
				return nil
			}
			if err != nil {
				glog.ErrorF("Kafka Consumer.FetchMessage error:%+v", err)
				continue
			}
			startTime := time.Now()
			err = handle(m)
			metricReqDuration.WithLabelValues(m.Topic, sub).Observe(float64(time.Since(startTime).Milliseconds()))
			result := "fail"
			if err == nil {
				result = "success"
				ackErr := kr.Reader.CommitMessages(ctx, m)
				if ackErr != nil {
					glog.ErrorF("Kafka Consumer.CommitMessages error:%+v", ackErr)
				}
			}
			metricsResult.WithLabelValues(m.Topic, pub, result).Inc()

			if err != nil {
				glog.ErrorF("Kafka Consumer.Handle with error:%+v", err)
			}
		}
	}
}

func (kr *Consumer) Close() error {
	kr.cancel()
	err := kr.Reader.Close()
	if err != nil {
		glog.ErrorF("Kafka Consumer close error:%v, conf:%#v", err, kr.Reader.Config())
	} else {
		glog.InfoF("Kafka Consumer close success, conf:%#v", kr.Reader.Config())
	}
	return err
}
