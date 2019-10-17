package nacos

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"io/ioutil"
)

type LocalNacos struct {
	FilePath string
	config_client.IConfigClient
}

func NewLocalNacos(filepath string) config_client.IConfigClient {
	return &LocalNacos{FilePath: filepath}
}

// 获取配置
func (l *LocalNacos) GetConfig(param vo.ConfigParam) (string, error) {
	// 判断必要参数
	if param.DataId != "" {
		return "", errors.New("The configuration file is incomplete.")
	}
	// 读取文件内容
	bytes, err := ioutil.ReadFile(l.FilePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 发布配置
func (l *LocalNacos) PublishConfig(param vo.ConfigParam) (bool, error) {
	return true, nil
}

// 删除配置
func (l *LocalNacos) DeleteConfig(param vo.ConfigParam) (bool, error) {
	return true, nil
}

// 监听配置
func (l *LocalNacos) ListenConfig(params vo.ConfigParam) (err error) {
	return nil
}
