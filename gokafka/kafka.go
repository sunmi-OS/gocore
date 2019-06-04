package gokafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"github.com/sunmi-OS/gocore/viper"
	"sync"
)

var producer *kafka.Writer


//construct the producer
func NewProducer()  {
	config := kafka.WriterConfig{
		Brokers:      viper.C.GetStringSlice("kafkaClient.brokers"),
		Topic:        viper.C.GetString("kafkaClient.topicName"),
		RequiredAcks: viper.C.GetInt("kafkaClient.acks"),
		Async:        viper.C.GetBool("kafkaClient.async"),
	}

	if viper.C.GetBool("kafkaClient.compression") == true {
		config.CompressionCodec =snappy.NewCompressionCodec() // snappy
	}

	producer = kafka.NewWriter(config)
}
//produce message
func Produce(msg []byte) error {
	return producer.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	})
}

//close the producer
func CloseProducer() error {
	if producer!=nil {
		return producer.Close()
	} else {
		return nil
	}
}


var consumer *kafka.Reader
var everyPartitionLastMessage sync.Map


//construct the consumer
func NewConsumer()  {
	consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   viper.C.GetStringSlice("kafkaClient.brokers"),
		GroupID:  viper.C.GetString("kafkaClient.consumerGroupId"),
		Topic:     viper.C.GetString("kafkaClient.topicName"),
		MinBytes:   viper.C.GetInt("kafkaClient.consumerMinBytes"),
		MaxBytes:  viper.C.GetInt("kafkaClient.consumerMaxBytes"),
	})
}

//consume message
func Consume() (kafka.Message,error) {
	lastMessage, err := consumer.FetchMessage(context.Background())
	if err == nil {
		everyPartitionLastMessage.Store(lastMessage.Partition, lastMessage)
	}
	return lastMessage,err
}

// commit offset for all partitions
func CommitOffsetForAllPartition(onCommit func(message kafka.Message) ) error {
	var err error
	background := context.Background()
	everyPartitionLastMessage.Range(func(key interface{}, value interface{}) bool {
		if err == nil {
			message := value.(kafka.Message)
			err = consumer.CommitMessages(background, message) 
			if err != nil {
				return false // stop iteration
			}
			everyPartitionLastMessage.Delete(key)
			if onCommit != nil {
				onCommit(message)
			}
			return true
		}
		return false
	})
	return err
}

//close the consumer
func CloseConsumer() error {
	if consumer != nil {
		return consumer.Close()
	} else {
		return nil
	}
}



