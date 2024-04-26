package config

import (
	"errors"
	"os"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/sunmi-OS/gocore/retry"
)

const (
	Runtime = "RUN_TIME"
	Local   = "local"

	_NamespaceId = "NAMESPACE_ID"
	_Endpoint    = "ENDPOINT"
	_AccessKey   = "ACCESS_KEY"
	_SecretKey   = "SECRET_KEY"

	_RegionId = "cn-hangzhou"
)

// 解析 acm 或本地 toml 配置到结构体
//
//	confPtr： 结构体指针
func ParseToml(group, dataId, localConfig string, confPtr interface{}) (err error) {
	beanValue := reflect.ValueOf(confPtr)
	if beanValue.Kind() != reflect.Ptr {
		return errors.New("confPtr must be ptr")
	}
	if beanValue.Elem().Kind() != reflect.Struct {
		return errors.New("confPtr must be struct ptr")
	}

	runtime := os.Getenv(Runtime)
	// local
	if runtime == Local && localConfig != "" {
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

	getConfigFunc := func() (err error) {
		client, err := clients.NewConfigClient(vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				TimeoutMs:   5000,
				NamespaceId: namespaceId,
				Endpoint:    endpoint,
				RegionId:    _RegionId,
				AccessKey:   accessKey,
				SecretKey:   secretKey,
				OpenKMS:     true,
				LogLevel:    "debug",
			},
		})
		if err != nil {
			return err
		}
		// 获取config信息
		config, err := client.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
		_, err = toml.Decode(config, confPtr)
		if err != nil {
			return err
		}
		return nil
	}
	return retry.Retry(getConfigFunc, 3, time.Second)
}

// 解析 acm 或本地 yaml 配置到结构体
//
//	confPtr： 结构体指针
func ParseYaml(group, dataId, localConfig string, confPtr interface{}) (err error) {
	beanValue := reflect.ValueOf(confPtr)
	if beanValue.Kind() != reflect.Ptr {
		return errors.New("confPtr must be ptr")
	}
	if beanValue.Elem().Kind() != reflect.Struct {
		return errors.New("confPtr must be struct ptr")
	}

	runtime := os.Getenv(Runtime)
	// local
	if runtime == Local && localConfig != "" {
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

	getConfigFunc := func() (err error) {
		client, err := clients.NewConfigClient(vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				TimeoutMs:   5000,
				NamespaceId: namespaceId,
				Endpoint:    endpoint,
				RegionId:    _RegionId,
				AccessKey:   accessKey,
				SecretKey:   secretKey,
				OpenKMS:     true,
				LogLevel:    "debug",
			},
		})
		if err != nil {
			return err
		}
		// 获取config信息
		config, err := client.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
		_, err = toml.Decode(config, confPtr)
		if err != nil {
			return err
		}
		return nil
	}
	return retry.Retry(getConfigFunc, 3, time.Second)
}
