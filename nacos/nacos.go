package nacos

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/viper"
)

type nacos struct {
	list              map[string]client
	runtime           string
	dataIdorGroupList []dataIdorGroup
	viperBase         string
	callbackList      map[string]func(namespace, group, dataId, data string)
	callbackRun       bool
	callbackFirst     sync.Map
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
	list:         make(map[string]client),
	callbackList: make(map[string]func(namespace, group, dataId, data string)),
	callbackRun:  false,
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

// 获取整套配置文件的拼接
func GetConfig() (string, error) {

	if nacosHarder.list[nacosHarder.runtime].localStrus {
		configs, err := nacosHarder.list[nacosHarder.runtime].cc.GetConfig(vo.ConfigParam{})
		if err != nil {
			return "", err
		}
		return configs, nil
	}

	var configs = ""

	for _, dataIdorGroup := range nacosHarder.dataIdorGroupList {
		group := dataIdorGroup.group
		for _, v := range dataIdorGroup.dataIds {

			content, err := nacosHarder.list[nacosHarder.runtime].cc.GetConfig(vo.ConfigParam{
				DataId: v,
				Group:  group})
			if err != nil {
				return "", err
			}
			configs += "\r\n" + content
			if nacosHarder.callbackRun == false {
				// 注册回调
				grouptodataId := group + v
				err := nacosHarder.list[nacosHarder.runtime].cc.ListenConfig(vo.ConfigParam{
					DataId:   v,
					Group:    group,
					OnChange: nacosHarder.callbackList[grouptodataId],
				})
				if err != nil {
					return "", err
				}
			}

		}
	}

	nacosHarder.callbackRun = true

	fmt.Println(configs)

	return configs, nil

}

// 设置环境变量
func SetRunTime(runtime string) {
	nacosHarder.runtime = runtime
}

// 设置需要读取哪些配置
func SetDataIds(group string, dataIds ...string) {
	nacosHarder.dataIdorGroupList = append(nacosHarder.dataIdorGroupList, dataIdorGroup{group: group, dataIds: dataIds})

	for _, v := range dataIds {
		grouptodataId := group + v
		nacosHarder.callbackList[grouptodataId] = func(namespace, group, dataId, data string) {
			i, _ := nacosHarder.callbackFirst.Load(grouptodataId)
			if cast.ToBool(i) == true {
				panic(errors.New(namespace + "\r\n" + group + "\r\n" + dataId + "\r\n" + data + "\r\n Updata Config"))
			}
			nacosHarder.callbackFirst.Store(grouptodataId, true)
		}
	}
}

// 配置回调方法
func SetCallBackFunc(group, dataId string, callbark func(namespace, group, dataId, data string)) {
	grouptodataId := group + dataId
	nacosHarder.callbackList[grouptodataId] = func(namespace, group, dataId, data string) {
		updateNacosToViper()
		i, _ := nacosHarder.callbackFirst.Load(grouptodataId)
		if cast.ToBool(i) == true {
			callbark(namespace, group, dataId, data)
		}
		nacosHarder.callbackFirst.Store(grouptodataId, true)
	}
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

func NacosToViper() {

	s, err := GetConfig()
	if err != nil {
		print(err)
	}
	viper.NewConfigToToml(s + nacosHarder.viperBase)
}

func SetviperBase(configs string) {
	nacosHarder.viperBase = configs
}

func NacosToViperFile(basefiles ...string) {

	if len(basefiles) > 0 {
		for _, v := range basefiles {
			bs, err := ioutil.ReadFile(v)
			if err != nil {
				panic(err)
			}
			nacosHarder.viperBase += "\r\n" + string(bs)
		}
	}

	s, err := GetConfig()
	if err != nil {
		print(err)
	}
	viper.NewConfigToToml(s + nacosHarder.viperBase)
}

func updateNacosToViper() {

	s, err := GetConfig()
	if err != nil {
		print(err)
	}
	viper.NewConfigToToml(s + nacosHarder.viperBase)
}
