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
		MinBytes:       10e3, // 10K
		MaxBytes:       10e6, // 10MB
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
		MinBytes:       10e3, // 10K
		MaxBytes:       10e6, // 10MB
		MaxWait:        time.Second,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	}
}

// NewConsumer conf每次重新生成
func NewConsumer(conf kafka.ReaderConfig) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Consumer{
		Reader: kafka.NewReader(conf),

		ctx:    ctx,
		cancel: cancel,
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

func (c *Consumer) Handle(ctx context.Context, handle func(msg kafka.Message) error) error {
	for {
		select {
		case <-ctx.Done():
			glog.WarnC(ctx, "Kafka Consumer ctx done, err=%+v", ctx.Err())
			return ctx.Err()
		case <-c.ctx.Done():
			glog.WarnC(ctx, "Kafka Consumer c.ctx done, err=%+v", c.ctx.Err())
			return c.ctx.Err()
		default:
			msg, err := c.Reader.ReadMessage(ctx)

			// io.EOF means consumer closed
			// io.ErrClosedPipe means committing messages on the consumer,
			// kafka will refire the messages on uncommitted messages, ignore
			if err == io.EOF || err == io.ErrClosedPipe {
				glog.WarnC(ctx, "Kafka Consumer ReadMessage failed, err=%+v(the reader has been closed)", err)
				return nil
			}
			if err != nil {
				glog.ErrorC(ctx, "Kafka Consumer ReadMessage failed, err=%+v", err)
				continue
			}

			startTime := time.Now()

			metricsDelay.WithLabelValues(msg.Topic).Observe(float64(time.Since(msg.Time).Milliseconds()))

			result := "success"
			if err = handle(msg); err != nil {
				result = "fail"
			}
			metricsResult.WithLabelValues(msg.Topic, sub, result).Inc()

			metricReqDuration.WithLabelValues(msg.Topic, sub).Observe(float64(time.Since(startTime).Milliseconds()))
		}
	}
}

func (c *Consumer) Close() error {
	c.cancel()

	err := c.Reader.Close()
	if err != nil {
		glog.ErrorF("Kafka Consumer close error:%v, conf:%#v", err, c.Reader.Config())
	} else {
		glog.InfoF("Kafka Consumer close success, conf:%#v", c.Reader.Config())
	}
	return err
}
