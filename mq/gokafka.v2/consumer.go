package gokafka

import (
	"context"
	"io"
	"time"

	"github.com/sunmi-OS/gocore/v2/conf/viper"

	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/segmentio/kafka-go"
)

/*
推荐命名：
后缀-S 代表消费(subscribe)
后缀-T 代表topic
后缀-P 代表生产(publish)

topic形式 {RunTime}-{产品线、描述}-T，如 Dev-ProductHelloMsg-T
消费group形式 {RunTime}-{AppID}-S，Dev-ProductCallbackJob-S
*/

// RunTimePrefix 带上大写开头环境前缀，例：ProductHello-T => Dev-ProductHello-T
func RunTimePrefix(str string) string {
	return utils.FirstUpper(utils.GetRunTime()) + "-" + str
}

// PreTopicPrefix Pre的topic，设置为Onl的preifx，例pre环境下：ProductHello-T => Onl-ProductHello-T
func PreTopicPrefix(str string) string {
	if utils.IsPre() {
		return utils.FirstUpper(utils.ReleaseEnv) + "-" + str
	}
	return utils.FirstUpper(utils.GetRunTime()) + "-" + str
}

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

func NewVipConsumerConfig(brokerKey string, groupIDKey string, topicKey string) kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:        viper.GetEnvConfig(brokerKey).SliceString(),
		GroupID:        viper.GetEnvConfig(groupIDKey).String(),
		Topic:          viper.GetEnvConfig(topicKey).String(),
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
	return &Consumer{
		ctx:    ctx,
		cancel: cancel,
		Reader: kafka.NewReader(conf),
	}
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
			metricReqDuration.WithLabelValues(m.Topic, "sub").Observe(float64(time.Since(startTime) / time.Millisecond))
			result := "fail"
			if err == nil {
				result = "success"
				ackErr := kr.Reader.CommitMessages(ctx, m)
				if ackErr != nil {
					glog.ErrorF("Kafka Consumer.CommitMessages error:%+v", ackErr)
				}
			}
			metricsResult.WithLabelValues(m.Topic, result).Inc()
		}
	}
}

// Close 别忘记调用该方法
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
