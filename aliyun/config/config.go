package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const (
	Runtime = "RUN_TIME"
	Local   = "local"
)

var (
	cc      config_client.IConfigClient
	runtime string
)

func Parse(group, dataId, localConfig string, confPtr interface{}) (err error) {
	runtime = os.Getenv(Runtime)

	// local
	if runtime == Local {
		if localConfig == "" {
			return errors.New("runtime is local, but local config is null")
		}
		_, err = toml.Decode(localConfig, confPtr)
		if err != nil {
			return err
		}
		return nil
	}

	// acm
	namespaceId := os.Getenv(_NamespaceId)
	endpoint := os.Getenv(_Endpoint)
	accessKey := os.Getenv(_AccessKey)
	secretKey := os.Getenv(_SecretKey)
	if endpoint == "" || namespaceId == "" || accessKey == "" || secretKey == "" {
		return errors.New("acm connection info error")
	}

	c, err := New(&Config{
		NamespaceId: namespaceId,
		Endpoint:    endpoint,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		RegionId:    _RegionId,
		LogLevel:    LogError,
	})
	if err != nil {
		return err
	}
	// 获取config信息
	config, err := c.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
	_, err = toml.Decode(config, confPtr)
	if err != nil {
		return err
	}
	cc = c
	return nil
}

func ListenConfig(group, dataId string, confPtr interface{}, f func(config string, err error)) error {
	if runtime != Local && cc != nil {
		// callback function
		callback := func(namespace string, group string, dataId string, data string) {
			_, err := toml.Decode(data, confPtr)
			if err != nil {
				f("", err)
				return
			}
			f(data, nil)
		}
		// 监听config
		cp := vo.ConfigParam{
			DataId:   dataId,
			Group:    group,
			OnChange: callback,
		}
		return cc.ListenConfig(cp)
	}
	return errors.New("runtime is local or nacos client is nil")
}
