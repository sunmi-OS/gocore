package nacos

import (
	"os"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type nacos struct {
	icc   config_client.IConfigClient
	vt    *ViperToml
	local bool
}

const (
	LogDebug = "debug"
	LogWarn  = "warn"
	LogError = "error"
	LogInfo  = "info"

	_NamespaceId     = "NAMESPACE_ID"
	_Endpoint        = "ENDPOINT"
	_AccessKey       = "ACCESS_KEY"
	_SecretKey       = "SECRET_KEY"
	_RegionId        = "REGION_ID"
	_DefaultRegionId = "cn-hangzhou"
)

var nacosHarder = &nacos{
	vt: &ViperToml{
		callbackList: make(map[string]func(namespace, group, dataId, data string)),
	},
}

// NewNacosEnv 注入ACM配置文件
// @TODO 需要兼容endpoint 和 service 两种方式
func NewNacosEnv() {
	namespaceId := os.Getenv(_NamespaceId)
	// 读取service地址，如果有service优先使用service连接方式
	endpoint := os.Getenv(_Endpoint)
	accessKey := os.Getenv(_AccessKey)
	secretKey := os.Getenv(_SecretKey)
	regionID := os.Getenv(_RegionId)
	if endpoint == "" || namespaceId == "" || accessKey == "" || secretKey == "" {
		panic("The configuration file cannot be empty.")
	}
	if regionID == "" {
		regionID = _DefaultRegionId
	}
	err := NewNacos(&constant.ClientConfig{
		Endpoint:    endpoint,
		NamespaceId: namespaceId,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		RegionId:    regionID,
		LogLevel:    LogError,
	})
	if err != nil {
		panic(err)
	}
}

// NewNacos 注入Nacos配置文件
func NewNacos(ccConfig *constant.ClientConfig, scConfigs ...constant.ServerConfig) error {
	defaultClientConfig(ccConfig)
	configClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  ccConfig,
		ServerConfigs: scConfigs,
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
func CallBackFunc(group, dataId string, callback func(namespace, group, dataId, data string)) error {
	err := nacosHarder.icc.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: callback,
	})
	if err != nil {
		return err
	}

	return nil
}
