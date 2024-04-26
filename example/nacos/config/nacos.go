package config

import (
	"os"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/sunmi-OS/gocore/nacos"
)

func InitNacos(runtime string) {
	nacos.SetRunTime(runtime)
	nacos.ViperTomlHarder.SetviperBase(baseConfig)
	switch runtime {
	case "local":
		nacos.AddLocalConfig(runtime, localConfig)
	default:
		Endpoint := os.Getenv("ENDPOINT")
		NamespaceId := os.Getenv("NAMESPACE_ID")
		AccessKey := os.Getenv("ACCESS_KEY")
		SecretKey := os.Getenv("SECRET_KEY")
		if Endpoint == "" || NamespaceId == "" || AccessKey == "" || SecretKey == "" {
			panic("The configuration file cannot be empty.")
		}
		err := nacos.AddAcmConfig(runtime, constant.ClientConfig{
			Endpoint:    Endpoint,
			NamespaceId: NamespaceId,
			AccessKey:   AccessKey,
			SecretKey:   SecretKey,
		})
		if err != nil {
			panic(err)
		}
	}
}
