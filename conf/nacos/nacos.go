package nacos

import (
	"io/ioutil"
	"os"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type nacos struct {
	icc   config_client.IConfigClient
	local bool
	vt    ViperToml
}

var nacosHarder = &nacos{
	vt: ViperToml{
		callbackList: make(map[string]func(namespace, group, dataId, data string)),
		callbackRun:  false,
	},
}

// SetLocalConfigFile 注入本地配置
func SetLocalConfigFile(filePath string) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	SetLocalConfig(string(bytes))
}

// SetLocalConfig 注入本地配置
func SetLocalConfig(configs string) {
	localNacos := NewLocalNacos(configs)
	nacosHarder.icc = localNacos
	nacosHarder.local = true
}

// NewNacos 注入Nacos配置文件
func NewNacos(ccConfig constant.ClientConfig, csConfigs ...constant.ServerConfig) error {
	defaultClientConfig(&ccConfig)
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": csConfigs,
		"clientConfig":  ccConfig,
	})
	if err != nil {
		return err
	}
	nacosHarder.icc = configClient
	return nil
}

// NewAcmEnv 注入ACM配置文件
func NewAcmEnv() {
	Endpoint := os.Getenv("ENDPOINT")
	NamespaceId := os.Getenv("NAMESPACE_ID")
	AccessKey := os.Getenv("ACCESS_KEY")
	SecretKey := os.Getenv("SECRET_KEY")

	if Endpoint == "" || NamespaceId == "" || AccessKey == "" || SecretKey == "" {
		panic("The configuration file cannot be empty.")
	}
	err := NewAcmConfig(&constant.ClientConfig{
		Endpoint:    Endpoint,
		NamespaceId: NamespaceId,
		AccessKey:   AccessKey,
		SecretKey:   SecretKey,
	})
	if err != nil {
		panic(err)
	}
}

// NewAcmConfig 注入ACM配置文件
func NewAcmConfig(ccConfig *constant.ClientConfig) error {
	defaultClientConfig(ccConfig)
	configClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  ccConfig,
		ServerConfigs: nil,
	})
	if err != nil {
		return err
	}
	nacosHarder.icc = configClient
	return nil
}

// defaultClientConfig 使用nacos时的默认配置
func defaultClientConfig(ccConfig *constant.ClientConfig) {
	if ccConfig.TimeoutMs == 0 {
		ccConfig.TimeoutMs = 5000
	}
	if ccConfig.BeatInterval == 0 {
		ccConfig.BeatInterval = 5 * 1000
	}
	if ccConfig.LogDir == "" {
		ccConfig.LogDir = "./nacos/logs"
	}
	if ccConfig.CacheDir == "" {
		ccConfig.CacheDir = "./nacos/cache"
	}
}

// GetConfig 获取单条配置
func GetConfig(group string, dataIds string) string {
	content, err := nacosHarder.icc.GetConfig(vo.ConfigParam{
		DataId: dataIds,
		Group:  group})
	if err != nil {
		return ""
	}
	return content
}

// CallBackFunc 参数更新回调方法
func CallBackFunc(group, dataId string, callbark func(namespace, group, dataId, data string)) error {
	err := nacosHarder.icc.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: callbark,
	})
	if err != nil {
		return err
	}

	return nil
}
