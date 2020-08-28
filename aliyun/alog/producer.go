package alog

import (
	"time"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
)

func Info(topic string, logs map[string]string) error {
	logs["level"] = "info"
	logmsg := producer.GenerateLog(uint32(time.Now().Unix()), logs)
	return logger.producer.SendLog(logger.Project, logger.LogStore, topic, logger.HostName, logmsg)
}

func Debug(topic string, logs map[string]string) error {
	logs["level"] = "debug"
	logmsg := producer.GenerateLog(uint32(time.Now().Unix()), logs)
	return logger.producer.SendLog(logger.Project, logger.LogStore, topic, logger.HostName, logmsg)
}

func Warn(topic string, logs map[string]string) error {
	logs["level"] = "warn"
	logmsg := producer.GenerateLog(uint32(time.Now().Unix()), logs)
	return logger.producer.SendLog(logger.Project, logger.LogStore, topic, logger.HostName, logmsg)
}

func Error(topic string, logs map[string]string) error {
	logs["level"] = "error"
	logmsg := producer.GenerateLog(uint32(time.Now().Unix()), logs)
	return logger.producer.SendLog(logger.Project, logger.LogStore, topic, logger.HostName, logmsg)
}
