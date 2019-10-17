package config

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/sunmi-OS/gocore/nacos"
)

var acmConfigs = map[string]constant.ClientConfig{
	"dev": constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "xxx",
		AccessKey:   "xxx",
		SecretKey:   "xxx",
	},
	"test": constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "xxx",
		AccessKey:   "xxx",
		SecretKey:   "xxx",
	},
	"uat": constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "xxx",
		AccessKey:   "xxx",
		SecretKey:   "xxx",
	},
	"onl": constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "xxx",
		AccessKey:   "xxx",
		SecretKey:   "xxx",
	},
}

func InitNacos(runtime string) {

	nacos.SetRunTime(runtime)

	nacos.SetviperBase(baseConfig)

	switch runtime {
	case "local":
		nacos.AddLocalConfig(runtime, localConfig)
	case "dev", "test", "uat", "onl":
		err := nacos.AddAcmConfig(runtime, acmConfigs[runtime])
		if err != nil {
			panic(err)
		}
	default:
		panic(errors.New("No corresponding configuration."))
	}

}
