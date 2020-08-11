package aliyunlog

import (
	"errors"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/sunmi-OS/gocore/viper"
	"os"
	"time"
)

type AliyunLog struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Project   string
	LogStore  string
	HostName  string
	Log       *producer.Producer
}

var LogClient AliyunLog

func InitLog(configName, LogStore string) {

	hostname, _ := os.Hostname()

	LogClient = AliyunLog{
		Project:   viper.GetEnvConfig(configName + ".Project"),
		Endpoint:  viper.GetEnvConfig(configName + ".Endpoint"),
		AccessKey: viper.GetEnvConfig(configName + ".AccessKey"),
		SecretKey: viper.GetEnvConfig(configName + ".SecretKey"),
		LogStore:  LogStore,
		HostName:  hostname,
	}
	err := checkConfig(LogClient)
	if err != nil {
		panic(err)
	}

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = LogClient.Endpoint
	producerConfig.AccessKeyID = LogClient.AccessKey
	producerConfig.AccessKeySecret = LogClient.SecretKey
	LogClient.Log = producer.InitProducer(producerConfig)

	LogClient.Log.Start()

	logmsg := producer.GenerateLog(uint32(time.Now().Unix()), map[string]string{"content": "log-start"})
	err = LogClient.Log.SendLog(LogClient.Project, LogClient.LogStore, "start", LogClient.HostName, logmsg)
	if err != nil {
		panic(err)
	}
}

func Info(topic string, logs map[string]string) error {

	logs["level"] = "info"
	logmsg := producer.GenerateLog(uint32(time.Now().Unix()), logs)
	return LogClient.Log.SendLog(LogClient.Project, LogClient.LogStore, topic, LogClient.HostName, logmsg)
}

func Error(topic string, logs map[string]string) error {

	logs["level"] = "error"
	logmsg := producer.GenerateLog(uint32(time.Now().Unix()), logs)
	return LogClient.Log.SendLog(LogClient.Project, LogClient.LogStore, topic, LogClient.HostName, logmsg)
}

func Close() {
	LogClient.Log.SafeClose()
}

func checkConfig(conf AliyunLog) (err error) {

	if conf.AccessKey == "" || conf.Endpoint == "" || conf.Project == "" || conf.LogStore == "" || conf.SecretKey == "" {
		err = errors.New("config Missing parameter")
	}
	return
}
