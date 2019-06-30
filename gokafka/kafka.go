package gokafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"github.com/sunmi-OS/gocore/viper"
	"sync"
	"time"
)



type Producer struct {
	producer *kafka.Writer
}

var producerMap sync.Map

func LoadProducerByTopic(topic string) *Producer {
	value, _ := producerMap.Load(topic)
	if value == nil {
		return nil
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
func (p *Producer) newProducer(topic string)  {


	viper.C.SetDefault(topic,map[string]interface{}{
		"acks":          -1,
		"async":          false,
		"compression":   true,
		"batchTimeout": 1,      //当你使用了 sync 方式写kafka，那么你需要将该值设置为1，因为不需要在进行无意义的等待时间
	})

	acks := viper.C.GetInt(topic + ".acks")
	async := viper.C.GetBool(topic + ".async")
	compression := viper.C.GetBool(topic + ".compression")
	batchTimeout := viper.C.GetInt(topic + ".batchTimeout")
	brokers := viper.C.GetStringSlice("kafkaClient.brokers")

	fmt.Printf("kafka producer config----->\n")
	fmt.Printf("topic: `%v`\n",topic)
	fmt.Printf("acks: `%v`\n",acks)
	fmt.Printf("async: `%v`\n",async)
	fmt.Printf("compression: `%v`\n",compression)
	fmt.Printf("batchTimeout: `%v`\n",batchTimeout)
	fmt.Printf("brokers: `%v`\n",brokers)
	fmt.Printf("kafka producer config<-----\n")

	config := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		RequiredAcks: acks,
		Async:        async,
		BatchTimeout: time.Millisecond * time.Duration(batchTimeout),
	}

	if compression  {
		config.CompressionCodec =snappy.NewCompressionCodec() // snappy
	}

	p.producer = kafka.NewWriter(config)
}

func (p *Producer) ProduceMsgs(msgs[]kafka.Message) error {
	return p.producer.WriteMessages(context.Background(), msgs...)
}

//produce message
func (p *Producer) ProduceWithKey(key []byte, value []byte) error {
	return p.producer.WriteMessages(context.Background(), kafka.Message{
		Key:   key,
		Value: value,
	})
}

//produce message
func (p *Producer) Produce(msg []byte) error {
	return p.producer.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	})
}

//close the producer
func (p *Producer) CloseProducer() error {
	if p.producer!=nil {
		return p.producer.Close()
	} else {
		return nil
	}
}


var consumer *kafka.Reader
var EveryPartitionLastMessage sync.Map


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
		EveryPartitionLastMessage.Store(lastMessage.Partition, lastMessage)
	}
	return lastMessage,err
}

// commit offset for all partitions
func CommitOffsetForAllPartition(onCommit func(message kafka.Message) ) error {
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

//close the consumer
func CloseConsumer() error {
	if consumer != nil {
		return consumer.Close()
	} else {
		return nil
	}
}



