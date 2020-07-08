package aliyunmq

import (
	rocketmq "github.com/apache/rocketmq-client-go/core"
	"sync"
)

type Consumer struct {
	RocketConf     RocketMQConfig
	ConsumerConfig *rocketmq.PushConsumerConfig
}

var ConsumerList sync.Map

func NewConsumer(configName string) (consumer Consumer) {

	conf := initConfig(configName)

	consumerConfig := &rocketmq.PushConsumerConfig{
		ClientConfig: rocketmq.ClientConfig{
			//您在阿里云 RocketMQ 控制台上申请的 GID。
			GroupID: conf.GroupID,
			//设置 TCP 协议接入点，从阿里云 RocketMQ 控制台的实例详情页面获取。 TCP方式接入VPN也不行
			NameServer: conf.NameServer,
			Credentials: &rocketmq.SessionCredentials{
				//您在阿里云账号管理控制台中创建的 AccessKeyId，用于身份认证。
				AccessKey: conf.AccessKey,
				//您在阿里云账号管理控制台中创建的 AccessKeySecret，用于身份认证。
				SecretKey: conf.SecretKey,
				//用户渠道，默认值为：ALIYUN。
				Channel: conf.Channel,
			},
		},
		//设置使用集群模式。
		Model: rocketmq.Clustering,
		//设置该消费者为普通消息消费。
		ConsumerModel: rocketmq.CoCurrently,
	}

	consumer = Consumer{
		RocketConf:     conf,
		ConsumerConfig: consumerConfig,
	}

	ConsumerList.LoadOrStore(configName, consumer)
	return
}

func (c Consumer) SetMessageModel(messageModel rocketmq.MessageModel) {
	c.ConsumerConfig.Model = messageModel
}

func (c Consumer) SetConsumerModel(consumerModel rocketmq.ConsumerModel) {
	c.ConsumerConfig.ConsumerModel = consumerModel
}

func GetPushConsumer(name string) (consumer rocketmq.PushConsumer, err error) {
	v, _ := ConsumerList.Load(name)
	conf := v.(Consumer)

	consumer, err = rocketmq.NewPushConsumer(conf.ConsumerConfig)
	if err != nil {
		return
	}
	return
}

