package rmqv1

import (
	"errors"
	"sync"

	rocketmq "github.com/apache/rocketmq-client-go/core"
	"github.com/sunmi-OS/gocore/viper"
)

var (
	ConsumerList sync.Map
	ProducerList sync.Map
)

func NewRocketMQ(configName string) (r *RocketMQ, err error) {
	r = &RocketMQ{
		ConfigName: configName,
		groupID:    viper.GetEnvConfig(configName + ".GroupID"),
		nameServer: viper.GetEnvConfig(configName + ".NameServer"),
		accessKey:  viper.GetEnvConfig(configName + ".AccessKey"),
		secretKey:  viper.GetEnvConfig(configName + ".SecretKey"),
		channel:    viper.GetEnvConfig(configName + ".Channel"),
	}

	if err = r.checkConfig(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RocketMQ) checkConfig() (err error) {
	if r.accessKey == "" {
		err = errors.New("missing AccessKey")
	}
	if r.secretKey == "" {
		err = errors.New("missing SecretKey")
	}
	if r.channel == "" {
		err = errors.New("missing Channel")
	}
	if r.groupID == "" {
		err = errors.New("missing GroupID")
	}
	if r.nameServer == "" {
		err = errors.New("missing NameServer")
	}
	return
}

func (r *RocketMQ) NewConsumer(c *ConsumerConfig) (consumer *Consumer, err error) {
	if c == nil {
		return nil, errors.New("ConsumerConfig can not be nil")
	}
	if c.Topic == "" {
		return nil, errors.New("topic can not be null")
	}
	pcc := &rocketmq.PushConsumerConfig{
		ClientConfig: rocketmq.ClientConfig{
			// 您在阿里云 RocketMQ 控制台上申请的 GID。
			GroupID: r.groupID,
			// 设置 TCP 协议接入点，从阿里云 RocketMQ 控制台的实例详情页面获取。 TCP方式接入VPN也不行
			NameServer: r.nameServer,
			Credentials: &rocketmq.SessionCredentials{
				// 您在阿里云账号管理控制台中创建的 AccessKeyId，用于身份认证。
				AccessKey: r.accessKey,
				// 您在阿里云账号管理控制台中创建的 AccessKeySecret，用于身份认证。
				SecretKey: r.secretKey,
				// 用户渠道，默认值为：ALIYUN。
				Channel: r.channel,
			},
		},
		// 设置使用集群模式。
		Model: rocketmq.Clustering,
		// 设置该消费者为普通消息消费。
		ConsumerModel: rocketmq.CoCurrently,
	}
	if c.MessageModel != 0 {
		pcc.Model = c.MessageModel
	}
	if c.ConsumerModel != 0 {
		pcc.ConsumerModel = c.ConsumerModel
	}
	pushConsumer, err := rocketmq.NewPushConsumer(pcc)
	if err != nil {
		return nil, err
	}
	ConsumerList.Store(r.ConfigName, pushConsumer)
	return &Consumer{pushConsumer: pushConsumer, cc: c}, nil
}

func (r *RocketMQ) NewProducer(c *ProducerConfig) (producer *Producer, err error) {
	if c == nil {
		return nil, errors.New("ProducerConfig can not be nil")
	}
	pc := &rocketmq.ProducerConfig{
		ClientConfig: rocketmq.ClientConfig{
			// 您在阿里云 RocketMQ 控制台上申请的 GID。
			GroupID: r.groupID,
			// 设置 TCP 协议接入点，从阿里云 RocketMQ 控制台的实例详情页面获取。 TCP方式接入VPN也不行
			NameServer: r.nameServer,
			Credentials: &rocketmq.SessionCredentials{
				// 您在阿里云账号管理控制台中创建的 AccessKeyId，用于身份认证。
				AccessKey: r.accessKey,
				// 您在阿里云账号管理控制台中创建的 AccessKeySecret，用于身份认证。
				SecretKey: r.secretKey,
				// 用户渠道，默认值为：ALIYUN。
				Channel: r.channel,
			},
		},
		// 主动设置该实例用于发送普通消息。
		ProducerModel: rocketmq.CommonProducer,
	}
	if c.ProducerModel != 0 {
		pc.ProducerModel = c.ProducerModel
	}
	p, err := rocketmq.NewProducer(pc)
	if err != nil {
		return nil, err
	}
	ProducerList.Store(r.ConfigName, p)
	return &Producer{producer: p}, nil
}
