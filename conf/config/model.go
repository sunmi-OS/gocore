package config

import "time"

const (
	LogDebug LogLevel = "debug"
	LogWarn  LogLevel = "warn"
	LogError LogLevel = "error"
	LogInfo  LogLevel = "info"

	_NamespaceId = "NAMESPACE_ID"
	_Endpoint    = "ENDPOINT"
	_AccessKey   = "ACCESS_KEY"
	_SecretKey   = "SECRET_KEY"

	_RegionId = "cn-hangzhou"
)

type LogLevel string

type Config struct {
	NamespaceId string
	Endpoint    string
	AccessKey   string
	SecretKey   string
	RegionId    string        // default: cn-hangzhou
	Timeout     time.Duration // default: 5s
	LogLevel    LogLevel      // default: warn
}
