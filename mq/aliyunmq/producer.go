package aliyunmq

import (
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var ProducerPool sync.Map

// NewProducer 初始化生产者
func NewProducer(configName string, option ...producer.Option) rocketmq.Producer {

	conf := initConfig(configName)

	producerOption := []producer.Option{
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.NameServer})),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: conf.AccessKey,
			SecretKey: conf.SecretKey,
		}),
		producer.WithNamespace(conf.Namespace),
	}

	if len(option) > 0 {
		producerOption = append(producerOption, option...)
	}

	conn, err := rocketmq.NewProducer(producerOption...)
	if err != nil {
		panic(err)
	}
	//请确保参数设置完成之后启动 Producer。
	err = conn.Start()
	if err != nil {
		panic(err)
	}

	ProducerPool.LoadOrStore(configName, conn)
	return conn
}

// GetProducer 获取生产者实例
func GetProducer(configName string) (conn rocketmq.Producer) {
	v, _ := ProducerPool.Load(configName)
	conn = v.(rocketmq.Producer)
	return
}

// CloseProducer 关闭所有生产者而连接
func CloseProducer() {
	ProducerPool.Range(func(key, value interface{}) bool {
		conn := value.(rocketmq.Producer)
		err := conn.Shutdown()
		return err == nil
	})
}
