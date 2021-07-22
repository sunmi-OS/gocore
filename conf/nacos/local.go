package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type LocalNacos struct {
	configs string
	config_client.IConfigClient
}

func NewLocalNacos(configs string) config_client.IConfigClient {
	return &LocalNacos{configs: configs}
}

func (l *LocalNacos) GetConfig(param vo.ConfigParam) (string, error) {
	str := l.configs
	return str, nil
}

func (l *LocalNacos) PublishConfig(param vo.ConfigParam) (bool, error) {
	return true, nil
}

func (l *LocalNacos) DeleteConfig(param vo.ConfigParam) (bool, error) {
	return true, nil
}

func (l *LocalNacos) ListenConfig(params vo.ConfigParam) (err error) {
	return nil
}
