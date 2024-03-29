package gokafka

import (
	"context"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
)

type Producer struct {
	producer *kafka.Writer
}

var producerMap sync.Map

func LoadProducerByTopic(topic string) *Producer {
	value, _ := producerMap.Load(topic)
	if value == nil {
		return Init(topic)
	} else {
		return value.(*Producer)
	}
}

func Init(topic string) *Producer {
	_, found := producerMap.Load(topic)
	if found {
		return nil //can not remap
	}
	producer := new(Producer)
	producer.newProducer(topic)
	producerMap.Store(topic, producer)
	return producer
}

//construct the producer
func (p *Producer) newProducer(topic string) {

	viper.C.SetDefault(topic, map[string]interface{}{
		"acks":         -1,
		"async":        false,
		"compression":  true,
		"batchTimeout": 1000,
	})

	acks := viper.GetEnvConfig(topic + ".acks").Int()
	async := viper.GetEnvConfig(topic + ".async").Bool()
	compression := viper.GetEnvConfig(topic + ".compression").Bool()
	batchTimeout := viper.GetEnvConfig(topic + ".batchTimeout").Int()
	brokers := viper.GetEnvConfig("kafkaClient.brokers").SliceString()

	// @TODO 更新kafka使用方式
	config := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		RequiredAcks: acks,
		Async:        async,
		BatchTimeout: time.Millisecond * time.Duration(batchTimeout),
	}

	if compression {
		config.CompressionCodec = snappy.NewCompressionCodec() // snappy
	}

	p.producer = kafka.NewWriter(config)
}

func (p *Producer) ProduceMsgs(msgs []kafka.Message) error {
	return p.producer.WriteMessages(context.Background(), msgs...)
}

// ProduceWithKey produce message
func (p *Producer) ProduceWithKey(key []byte, value []byte) error {
	return p.producer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})
}

// Produce produce message
func (p *Producer) Produce(msg []byte) error {
	return p.producer.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	})
}

// CloseProducer close the producer
func (p *Producer) CloseProducer() error {
	if p.producer != nil {
		return p.producer.Close()
	} else {
		return nil
	}
}

var consumer *kafka.Reader
var EveryPartitionLastMessage sync.Map

// NewConsumer construct the consumer
func NewConsumer() {
	consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  viper.GetEnvConfig("kafkaClient.brokers").SliceString(),
		GroupID:  viper.GetEnvConfig("kafkaClient.consumerGroupId").String(),
		Topic:    viper.GetEnvConfig("kafkaClient.topicName").String(),
		MinBytes: viper.GetEnvConfig("kafkaClient.consumerMinBytes").Int(),
		MaxBytes: viper.GetEnvConfig("kafkaClient.consumerMaxBytes").Int(),
	})
}

// Consume consume message
func Consume() (kafka.Message, error) {
	lastMessage, err := consumer.FetchMessage(context.Background())
	if err == nil {
		EveryPartitionLastMessage.Store(lastMessage.Partition, lastMessage)
	}
	return lastMessage, err
}

func ConsumeByCallback(consume func(kafka.Message, error) bool) {
	lastMessage, err := consumer.FetchMessage(context.Background())
	goon := consume(lastMessage, err)
	if goon {
		EveryPartitionLastMessage.Store(lastMessage.Partition, lastMessage)
	}
}

// CommitOffsetForAllPartition commit offset for all partitions
func CommitOffsetForAllPartition(onCommit func(message kafka.Message)) error {
	var err error
	background := context.Background()
	EveryPartitionLastMessage.Range(func(key interface{}, value interface{}) bool {
		if err == nil {
			message := value.(kafka.Message)
			err = consumer.CommitMessages(background, message)
			if err != nil {
				return false // stop iteration
			}
			EveryPartitionLastMessage.Delete(key)
			if onCommit != nil {
				onCommit(message)
			}
			return true
		}
		return false
	})
	return err
}

// CloseConsumer close the consumer
func CloseConsumer() error {
	if consumer != nil {
		return consumer.Close()
	} else {
		return nil
	}
}
