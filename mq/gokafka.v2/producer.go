package gokafka

import (
	"context"
	"time"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
	ctx    context.Context
	cancel context.CancelFunc
}

func NewProducerConfig(brokers []string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Async: true,
	}
}

func NewVipProducerConfig(brokerKey string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(viper.GetEnvConfig(brokerKey).SliceString()...),
		Async: true,
	}
}

// NewProducer conf每次重新生成
func NewProducer(conf *kafka.Writer) *Producer {
	glog.InfoF("start one kafka producer, conf:%#v", conf)
	ctx, cancel := context.WithCancel(context.Background())
	return &Producer{
		ctx:    ctx,
		cancel: cancel,
		Writer: conf,
	}
}

func (w *Producer) Send(ctx context.Context, topic string, key string, value []byte) error {
	startTime := time.Now()
	err := w.Writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
	result := "success"
	if err != nil {
		result = "fail"
		glog.ErrorF("Kafka WriteMessages unexpected error:%v", err)
	}
	metricsResult.WithLabelValues(topic, result).Inc()
	metricReqDuration.WithLabelValues(w.Writer.Topic, "pub").Observe(float64(time.Since(startTime) / time.Millisecond))
	return err
}

func (w *Producer) SendBatch(ctx context.Context, msgs ...kafka.Message) error {
	startTime := time.Now()
	err := w.Writer.WriteMessages(ctx, msgs...)
	result := "success"
	if err != nil {
		result = "fail"
		glog.ErrorF("Kafka WriteMessages unexpected error:%v", err)
	}
	cost := float64(time.Since(startTime) / time.Millisecond)
	for _, msg := range msgs {
		metricsResult.WithLabelValues(msg.Topic, result).Inc()
		metricReqDuration.WithLabelValues(msg.Topic, "pub").Observe(cost)
	}
	return err
}

func (w *Producer) Close() error {
	w.cancel()
	err := w.Writer.Close()
	if err != nil {
		glog.ErrorF("Kafka Producer close error:%v, conf:%#v", err, w.Writer)
	} else {
		glog.InfoF("Kafka Producer close success, conf:%#v", w.Writer)
	}
	return err
}
