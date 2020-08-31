package mq

import (
	"sync"

	rocketmq "github.com/apache/rocketmq-client-go/core"
)

type Producer struct {
	RocketConf     *RocketMQConfig
	ProducerConfig *rocketmq.ProducerConfig
}

var ProducerList sync.Map

func NewProducer(configName string) (producer Producer) {

	conf := initConfig(configName)

	producer = Producer{
		RocketConf: conf,
		ProducerConfig: &rocketmq.ProducerConfig{
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
			//主动设置该实例用于发送普通消息。
			ProducerModel: rocketmq.CommonProducer,
		},
	}

	conn, err := rocketmq.NewProducer(producer.ProducerConfig)
	if err != nil {
		panic(err)
	}
	//请确保参数设置完成之后启动 Producer。
	err = conn.Start()
	if err != nil {
		panic(err)
	}

	ProducerList.Store(configName, conn)
	return
}

func (p Producer) SetProducerModel(producerModel rocketmq.ProducerModel) {
	p.ProducerConfig.ProducerModel = producerModel
}

func GetProducer(name string) (producer rocketmq.Producer) {
	v, _ := ProducerList.Load(name)
	producer = v.(rocketmq.Producer)
	return
}
