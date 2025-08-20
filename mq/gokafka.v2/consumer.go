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

	autoCommit bool // 是否自动提交offset

	ctx    context.Context
	cancel context.CancelFunc
}

type Option func(*Consumer)

func AutoCommit() Option {
	return func(c *Consumer) {
		c.autoCommit = true
	}
}

func NewConsumerConfig(brokers []string, groupID string, topic string) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		MinBytes:    10e3, // 10K
		MaxBytes:    10e6, // 10MB
		MaxWait:     time.Second,
		StartOffset: kafka.LastOffset,
	}
}

func NewVipConsumerConfig(brokername string, groupID string, topic string) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:     viper.GetEnvConfig(brokername + ".Brokers").SliceString(),
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    10e3, // 10K
		MaxBytes:    10e6, // 10MB
		MaxWait:     time.Second,
		StartOffset: kafka.LastOffset,
	}
}

// NewConsumer conf每次重新生成
func NewConsumer(conf kafka.ReaderConfig, opts ...Option) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Consumer{
		ctx:    ctx,
		cancel: cancel,
	}

	for _, opt := range opts {
		opt(c)
	}
	if c.autoCommit {
		conf.CommitInterval = time.Second // 自动提交offset
	}
	c.Reader = kafka.NewReader(conf)

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
	if kr.autoCommit {
		return kr.handleWithAutoCommit(ctx, handle)
	}
	return kr.handleWithManualCommit(ctx, handle)
}

func (kr *Consumer) handleWithAutoCommit(ctx context.Context, handle func(msg kafka.Message) error) error {
	for {
		select {
		case <-ctx.Done():
			glog.WarnC(ctx, "Kafka Consumer(AutoCommit) ctx done, err=%+v", ctx.Err())
			return ctx.Err()
		case <-kr.ctx.Done():
			glog.WarnC(ctx, "Kafka Consumer(AutoCommit) kr.ctx done, err=%+v", kr.ctx.Err())
			return kr.ctx.Err()
		default:
			msg, err := kr.Reader.FetchMessage(ctx)

			// io.EOF means consumer closed
			// io.ErrClosedPipe means committing messages on the consumer,
			// kafka will refire the messages on uncommitted messages, ignore
			if err == io.EOF || err == io.ErrClosedPipe {
				glog.WarnC(ctx, "Kafka Consumer(AutoCommit) FetchMessage failed, err=%+v(the reader has been closed)", err)
				return nil
			}
			if err != nil {
				glog.ErrorC(ctx, "Kafka Consumer(AutoCommit) FetchMessage failed, err=%+v", err)
				continue
			}

			startTime := time.Now()

			result := "success"
			if err = handle(msg); err != nil {
				result = "fail"
			}
			metricsResult.WithLabelValues(msg.Topic, sub, result).Inc()

			metricReqDuration.WithLabelValues(msg.Topic, sub).Observe(float64(time.Since(startTime).Milliseconds()))
			metricsDelay.WithLabelValues(msg.Topic).Observe(float64(time.Since(msg.Time).Milliseconds()))
		}
	}
}

func (kr *Consumer) handleWithManualCommit(ctx context.Context, handle func(msg kafka.Message) error) error {
	for {
		select {
		case <-ctx.Done():
			glog.WarnC(ctx, "Kafka Consumer(ManualCommit) ctx done, err=%+v", ctx.Err())
			return ctx.Err()
		case <-kr.ctx.Done():
			glog.WarnC(ctx, "Kafka Consumer(ManualCommit) kr.ctx done, err=%+v", kr.ctx.Err())
			return kr.ctx.Err()
		default:
			msg, err := kr.Reader.FetchMessage(ctx)

			// io.EOF means consumer closed
			// io.ErrClosedPipe means committing messages on the consumer,
			// kafka will refire the messages on uncommitted messages, ignore
			if err == io.EOF || err == io.ErrClosedPipe {
				glog.WarnC(ctx, "Kafka Consumer(ManualCommit) FetchMessage failed, err=%+v(the reader has been closed)", err)
				return nil
			}
			if err != nil {
				glog.ErrorC(ctx, "Kafka Consumer(ManualCommit) FetchMessage failed, err=%+v", err)
				continue
			}

			startTime := time.Now()

			err = handle(msg)

			metricReqDuration.WithLabelValues(msg.Topic, sub).Observe(float64(time.Since(startTime).Milliseconds()))
			metricsDelay.WithLabelValues(msg.Topic).Observe(float64(time.Since(msg.Time).Milliseconds()))

			if ackErr := kr.Reader.CommitMessages(ctx, msg); ackErr != nil {
				glog.ErrorC(ctx, "Kafka Consumer(ManualCommit) CommitMessages failed, err=%+v", ackErr)
			}

			if err != nil {
				metricsResult.WithLabelValues(msg.Topic, sub, "fail").Inc()
				continue
			}
			metricsResult.WithLabelValues(msg.Topic, sub, "success").Inc()
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
