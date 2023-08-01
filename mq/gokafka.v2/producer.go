package gokafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils/closes"

	"github.com/segmentio/kafka-go"
)

var ProducerPool sync.Map
var closeOnce sync.Once

type Producer struct {
	Writer     *kafka.Writer
	ctx        context.Context
	cancel     context.CancelFunc
	configName string
}

// NewProducerConfig 该方法返回值不能复用，每次NewProducer时都需要调用一次
func NewProducerConfig(brokers []string) *kafka.Writer {
	if len(brokers) == 0 {
		glog.Error("Kafka Brokers is empty")
		return nil
	}
	return &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Async: true,
	}
}

// NewVipProducerConfig 该方法返回值不能复用，每次NewProducer时都需要调用一次
func NewVipProducerConfig(configName string) *kafka.Writer {
	brokers := viper.GetEnvConfig(configName + ".Brokers").SliceString()
	if len(brokers) == 0 {
		glog.ErrorF("Kafka configName:%v Brokers is empty", configName)
		return nil
	}
	return &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Async: true,
	}
}

// NewProducer conf每次重新生成
func NewProducer(configName string, conf *kafka.Writer) *Producer {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Producer{
		ctx:        ctx,
		cancel:     cancel,
		Writer:     conf,
		configName: configName,
	}

	oldProducer, _ := ProducerPool.Load(configName)
	ProducerPool.Store(configName, p)
	if oldProducer != nil {
		if pd, _ := oldProducer.(*Producer); pd != nil {
			glog.InfoF("Kafka has same configName producer, close it, configName:%v", configName)
			pd.cancel()
			pd.Writer.Close()
		}
	}

	closeOnce.Do(func() {
		closes.AddShutdown(closes.ModuleClose{
			Name:     "Kafka Producer Close",
			Priority: closes.MQPriority,
			Func:     Close,
		})
	})
	glog.InfoF("Kafka start one conf= %v", formatWriterConfig(configName, conf))
	return p
}

func formatWriterConfig(configName string, conf *kafka.Writer) string {
	return fmt.Sprintf("configName:%v Addr:%v Async:%v", configName, conf.Addr, conf.Async)
}

func GetProducer(configName string) (conn *Producer) {
	if conn, ok := ProducerPool.Load(configName); ok {
		return conn.(*Producer)
	}
	return nil
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
		glog.ErrorC(ctx, "Kafka WriteMessages unexpected error:%v", err)
	}
	metricsResult.WithLabelValues(topic, pub, result).Inc()
	metricReqDuration.WithLabelValues(topic, pub).Observe(float64(time.Since(startTime).Milliseconds()))
	return err
}

func (w *Producer) SendBatch(ctx context.Context, msgs ...kafka.Message) error {
	startTime := time.Now()
	err := w.Writer.WriteMessages(ctx, msgs...)
	result := "success"
	if err != nil {
		result = "fail"
		glog.ErrorC(ctx, "Kafka WriteMessages unexpected error:%v", err)
	}
	cost := float64(time.Since(startTime).Milliseconds())
	for _, msg := range msgs {
		metricsResult.WithLabelValues(msg.Topic, pub, result).Inc()
		metricReqDuration.WithLabelValues(msg.Topic, pub).Observe(cost)
	}
	return err
}

func Close() {
	ProducerPool.Range(func(key, value interface{}) bool {
		glog.InfoF("Kafka Producer start close key: %s", key)
		ProducerPool.Delete(key)
		if p, _ := value.(*Producer); p != nil {
			p.cancel()
			err := p.Writer.Close()
			if err != nil {
				glog.ErrorF("Kafka Producer close error:%v, conf= %v", err, formatWriterConfig(p.configName, p.Writer))
			} else {
				glog.InfoF("Kafka Producer close success, conf= %v", formatWriterConfig(p.configName, p.Writer))
			}
		}
		return true
	})
}
