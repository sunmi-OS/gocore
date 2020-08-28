package alog

import (
	"os"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/sunmi-OS/gocore/ecode"
	"github.com/sunmi-OS/gocore/viper"
)

type Logger struct {
	Project  string
	LogStore string
	HostName string
	// private
	endpoint  string
	accessKey string
	secretKey string
	producer  *producer.Producer
}

var logger *Logger

func New(c *LoggerConfig) {
	hostname, _ := os.Hostname()
	logger = &Logger{
		Project:   viper.GetEnvConfig(c.ConfigName + ".Project"),
		LogStore:  c.LogStore,
		HostName:  hostname,
		endpoint:  viper.GetEnvConfig(c.ConfigName + ".Endpoint"),
		accessKey: viper.GetEnvConfig(c.ConfigName + ".AccessKey"),
		secretKey: viper.GetEnvConfig(c.ConfigName + ".SecretKey"),
		producer:  nil,
	}
	if err := logger.checkConfig(); err != nil {
		panic(err)
	}
	logger.newProducer().Start()
}

func (l *Logger) newProducer() *Logger {
	pc := producer.GetDefaultProducerConfig()
	pc.Endpoint = l.endpoint
	pc.AccessKeyID = l.accessKey
	pc.AccessKeySecret = l.secretKey
	l.producer = producer.InitProducer(pc)
	return l
}

func (l *Logger) Start() {
	l.producer.Start()
}

func (l *Logger) checkConfig() error {
	if l.accessKey == "" || l.endpoint == "" || l.Project == "" || l.LogStore == "" || l.secretKey == "" {
		return ecode.ConfigParamErr
	}
	return nil
}
