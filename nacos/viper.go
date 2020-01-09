package nacos

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/viper"
	"io/ioutil"
	"sync"
	"time"
)

type ViperToml struct {
	dataIdorGroupList []dataIdorGroup
	viperBase         string
	callbackList      map[string]func(namespace, group, dataId, data string)
	callbackRun       bool
	callbackFirst     sync.Map
}

var ViperTomlHarder = &ViperToml{
	callbackList: make(map[string]func(namespace, group, dataId, data string)),
	callbackRun:  false,
}

// 获取整套配置文件的拼接
func (vt *ViperToml) GetConfig() (string, error) {

	if nacosHarder.list[nacosHarder.runtime].localStrus {
		configs, err := nacosHarder.list[nacosHarder.runtime].cc.GetConfig(vo.ConfigParam{})
		if err != nil {
			return "", err
		}
		return configs, nil
	}

	var configs = ""

	for _, dataIdorGroup := range vt.dataIdorGroupList {
		group := dataIdorGroup.group
		for _, v := range dataIdorGroup.dataIds {

			content, err := nacosHarder.list[nacosHarder.runtime].cc.GetConfig(vo.ConfigParam{
				DataId: v,
				Group:  group})
			if err != nil {
				return "", err
			}
			configs += "\r\n" + content
			if vt.callbackRun == false {
				// 注册回调
				grouptodataId := group + v
				err := nacosHarder.list[nacosHarder.runtime].cc.ListenConfig(vo.ConfigParam{
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

	vt.callbackRun = true

	return configs, nil

}

// 设置需要读取哪些配置
func (vt *ViperToml) SetDataIds(group string, dataIds ...string) {
	vt.dataIdorGroupList = append(vt.dataIdorGroupList, dataIdorGroup{group: group, dataIds: dataIds})

	for _, v := range dataIds {
		grouptodataId := group + v
		vt.callbackList[grouptodataId] = func(namespace, group, dataId, data string) {
			i, _ := vt.callbackFirst.Load(grouptodataId)
			if cast.ToBool(i) == true {
				panic(errors.New(namespace + "\r\n" + group + "\r\n" + dataId + "\r\n" + data + "\r\n Updata Config"))
			}
			vt.callbackFirst.Store(grouptodataId, true)
		}
	}
}

// 配置回调方法
func (vt *ViperToml) SetCallBackFunc(group, dataId string, callbark func(namespace, group, dataId, data string)) {
	grouptodataId := group + dataId
	vt.callbackList[grouptodataId] = func(namespace, group, dataId, data string) {
		vt.updateNacosToViper()
		i, _ := vt.callbackFirst.Load(grouptodataId)
		if cast.ToBool(i) == true {
			callbark(namespace, group, dataId, data)
		}
		vt.callbackFirst.Store(grouptodataId, true)
	}
}

func (vt *ViperToml) NacosToViper() {

	var err error
	s := ""

	for i := 0; i < 3; i++ {
		s, err = vt.GetConfig()
		if err != nil {
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		panic(err)
	}

	viper.NewConfigToToml(s + vt.viperBase)
}

//
func (vt *ViperToml) SetviperBase(configs string) {
	vt.viperBase = configs
}

//
func (vt *ViperToml) NacosToViperFile(basefiles ...string) {

	if len(basefiles) > 0 {
		for _, v := range basefiles {
			bs, err := ioutil.ReadFile(v)
			if err != nil {
				panic(err)
			}
			vt.viperBase += "\r\n" + string(bs)
		}
	}

	s, err := vt.GetConfig()
	if err != nil {
		print(err)
	}
	viper.NewConfigToToml(s + vt.viperBase)
}

func (vt *ViperToml) updateNacosToViper() {

	s, err := vt.GetConfig()
	if err != nil {
		print(err)
	}
	viper.MerageConfigToToml(s + vt.viperBase)
}
