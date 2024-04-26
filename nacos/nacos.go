package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type nacos struct {
	list    map[string]client
	runtime string
	vt      ViperToml
}

type client struct {
	cc         config_client.IConfigClient
	localStrus bool
}

type dataIdorGroup struct {
	group   string
	dataIds []string
}

var nacosHarder = &nacos{
	list: make(map[string]client),
	vt: ViperToml{
		callbackList: make(map[string]func(namespace, group, dataId, data string)),
		callbackRun:  false,
	},
}

// 注入本地文件配置
func AddLocalConfigFile(runTime string, filePath string) {
	localNacos := NewLocalNacos(filePath)
	nacosHarder.list[runTime] = client{cc: localNacos, localStrus: true}
}

// 注入本地文件配置
func AddLocalConfig(runTime string, configs string) {
	localNacos := NewLocalNacos(configs)
	nacosHarder.list[runTime] = client{cc: localNacos, localStrus: true}
}

// 注入Nacos配置文件
func AddNacosConfig(runTime string, ccConfig constant.ClientConfig, csConfigs []constant.ServerConfig) error {

	defaltClientConfig(&ccConfig)

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": csConfigs,
		"clientConfig":  ccConfig,
	})
	if err != nil {
		return err
	}
	nacosHarder.list[runTime] = client{cc: configClient}
	return nil
}

// 注入ACM配置文件
func AddAcmConfig(runTime string, ccConfig constant.ClientConfig) error {

	defaltClientConfig(&ccConfig)

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": ccConfig,
	})
	if err != nil {
		return err
	}
	nacosHarder.list[runTime] = client{cc: configClient}

	return nil
}

// 设置环境变量
func SetRunTime(runtime string) {
	nacosHarder.runtime = runtime
}

// 使用nacos时的默认配置
func defaltClientConfig(ccConfig *constant.ClientConfig) {
	if ccConfig.TimeoutMs == 0 {
		ccConfig.TimeoutMs = 5 * 1000
	}
	if ccConfig.ListenInterval == 0 {
		ccConfig.ListenInterval = 30 * 1000
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

func GetConfig(group string, dataIds string) string {
	content, err := nacosHarder.list[nacosHarder.runtime].cc.GetConfig(vo.ConfigParam{
		DataId: dataIds,
		Group:  group})
	if err != nil {
		return ""
	}
	return content
}

// 配置回调方法
func CallBackFunc(group, dataId string, callbark func(namespace, group, dataId, data string)) error {

	err := nacosHarder.list[nacosHarder.runtime].cc.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: callbark,
	})
	if err != nil {
		return err
	}

	return nil
}
