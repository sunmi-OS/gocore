package config

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/sunmi-OS/gocore/nacos"
)

var acmConfigs = map[string]constant.ClientConfig{
	"dev": constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "96559b31-4749-4510-8425-2209ebfd9155",
		AccessKey:   "LTAI4Fvzja6U3zBHcUcecoNz",
		SecretKey:   "aJdP1cRJYnSrV9KJKlbhTUZseLlWRb",
	},
	"test": constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "96559b31-4749-4510-8425-2209ebfd9155",
		AccessKey:   "LTAI4Fvzja6U3zBHcUcecoNz",
		SecretKey:   "aJdP1cRJYnSrV9KJKlbhTUZseLlWRb",
	},
	"uat": constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "96559b31-4749-4510-8425-2209ebfd9155",
		AccessKey:   "LTAI4Fvzja6U3zBHcUcecoNz",
		SecretKey:   "aJdP1cRJYnSrV9KJKlbhTUZseLlWRb",
	},
	"onl": constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "96559b31-4749-4510-8425-2209ebfd9155",
		AccessKey:   "LTAI4Fvzja6U3zBHcUcecoNz",
		SecretKey:   "aJdP1cRJYnSrV9KJKlbhTUZseLlWRb",
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
