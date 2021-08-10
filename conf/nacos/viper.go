package nacos

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/sunmi-OS/gocore/v2/conf/viper"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

type configParam struct {
	group   string
	dataIds []string
}

type ViperToml struct {
	dataIdOrGroupList []configParam
	callbackList      map[string]func(namespace, group, dataId, data string)
}

// GetViper 获取VT实例
func GetViper() *ViperToml {
	return nacosHarder.vt
}

// GetConfig 获取整套配置文件
func (vt *ViperToml) GetConfig() (string, error) {
	if nacosHarder.local {
		configs, err := nacosHarder.icc.GetConfig(vo.ConfigParam{})
		if err != nil {
			return "", err
		}
		return configs, nil
	}
	var configs = ""
	for _, dataIdorGroup := range vt.dataIdOrGroupList {
		group := dataIdorGroup.group
		for _, v := range dataIdorGroup.dataIds {
			content, err := nacosHarder.icc.GetConfig(vo.ConfigParam{
				DataId: v,
				Group:  group})
			if err != nil {
				return "", err
			}
			configs += "\r\n" + content
			if len(vt.callbackList) > 0 {
				// 注册回调
				grouptodataId := group + v
				err := nacosHarder.icc.ListenConfig(vo.ConfigParam{
					DataId:   v,
					Group:    group,
					OnChange: vt.callbackList[grouptodataId],
				})
				if err != nil {
					return "", err
				}
			}
		}
	}
	return configs, nil
}

// GetConfigParse 获取配置并且绑定结构体
func (vt *ViperToml) GetConfigParse(confPtr interface{}) error {
	config, err := vt.GetConfig()
	if err != nil {
		return err
	}
	_, err = toml.Decode(config, confPtr)
	if err != nil {
		return err
	}
	return nil
}

// SetDataIds 设置需要读取哪些配置
func (vt *ViperToml) SetDataIds(group string, dataIds ...string) {
	vt.dataIdOrGroupList = append(vt.dataIdOrGroupList, configParam{group: group, dataIds: dataIds})
	for _, v := range dataIds {
		groupToDataId := group + v
		vt.callbackList[groupToDataId] = func(namespace, group, dataId, data string) {
			panic(errors.New(namespace + "\r\n" + group + "\r\n" + dataId + "\r\n" + data + "\r\n Updata Config"))
		}
	}
}

// SetCallBackFunc 配置回调方法
func (vt *ViperToml) SetCallBackFunc(group, dataId string, callback func(namespace, group, dataId, data string)) {
	groupToDataId := group + dataId
	vt.callbackList[groupToDataId] = func(namespace, group, dataId, data string) {
		vt.NacosToViper()
		callback(namespace, group, dataId, data)
	}
}

// NacosToViper 同步Nacos读取的配置注入Viper
func (vt *ViperToml) NacosToViper() {
	s, err := vt.GetConfig()
	if err != nil {
		panic(err)
	}
	viper.MergeConfigToToml(s)
}

// SetViperBase 注入基础配置
func (vt *ViperToml) SetViperBase(configs string) {
	viper.MergeConfigToToml(configs)
}
