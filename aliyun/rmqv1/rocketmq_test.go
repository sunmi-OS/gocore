package rmqv1

import (
	"log"
	"testing"

	rocketmq "github.com/apache/rocketmq-client-go/core"
	"github.com/sunmi-OS/gocore/viper"
)

func TestConsumer(t *testing.T) {
	viper.C.SetDefault("aliyunmq.GroupID", "GID_xxxx")
	viper.C.SetDefault("aliyunmq.NameServer", "http://xxxx.cn-hangzhou.mq-internal.aliyuncs.com:8080")
	viper.C.SetDefault("aliyunmq.AccessKey", "xxxx")
	viper.C.SetDefault("aliyunmq.SecretKey", "xxxx")
	viper.C.SetDefault("aliyunmq.Channel", "ALIYUN")

	rocketMQ, err := NewRocketMQ("aliyunmq")
	if err != nil {
		log.Println(err)
		return
	}

	cc := &ConsumerConfig{
		Topic: "sunmi",
	}

	consumer, err := rocketMQ.NewConsumer(cc)
	if err != nil {
		log.Println("rocketMQ.NewConsumer err:", err)
		return
	}
	err = consumer.Subscribe(func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {
		// todo: do something

		return rocketmq.ConsumeSuccess
	})
	if err != nil {
		log.Println("consumer.Subscribe err:", err)
		return
	}
	if err = consumer.Start(); err != nil {
		log.Println("consumer.Start err:", err)
		return
	}
	defer consumer.Close()
	ch := make(chan interface{})
	<-ch
}

func TestProducer(t *testing.T) {
	viper.C.SetDefault("aliyunmq.GroupID", "GID_xxxx")
	viper.C.SetDefault("aliyunmq.NameServer", "http://xxxx.cn-hangzhou.mq-internal.aliyuncs.com:8080")
	viper.C.SetDefault("aliyunmq.AccessKey", "xxxx")
	viper.C.SetDefault("aliyunmq.SecretKey", "xxxx")
	viper.C.SetDefault("aliyunmq.Channel", "ALIYUN")

	rocketMQ, err := NewRocketMQ("aliyunmq")
	if err != nil {
		log.Println(err)
		return
	}

	producer, err := rocketMQ.NewProducer(&ProducerConfig{})
	if err != nil {
		log.Println("rocketMQ.NewProducer err:", err)
		return
	}

	if err = producer.Start(); err != nil {
		log.Println("producer.Start err:", err)
		return
	}

	result, err := producer.SendMessageSync(&rocketmq.Message{Keys: "orderid", Topic: "xxxxxxx", Body: "message"})
	if err != nil {
		log.Println("producer.SendMessageSync err:", err)
		return
	}
	log.Println(result)
	producer.Close()
}
