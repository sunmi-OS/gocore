package config

import (
	"time"

	nacos "github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/sunmi-OS/gocore/utils/retry"
)

func New(cfg *Config) (nacosCli config_client.IConfigClient, err error) {

	if cfg.Timeout == 0 {
		cfg.Timeout = time.Second * 5
	}
	if cfg.RegionId == "" {
		cfg.RegionId = _RegionId
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = LogWarn
	}

	initNacos := func() error {
		nacosCli, err = nacos.NewConfigClient(vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				TimeoutMs:   uint64(cfg.Timeout.Milliseconds()),
				NamespaceId: cfg.NamespaceId,
				Endpoint:    cfg.Endpoint,
				RegionId:    cfg.RegionId,
				AccessKey:   cfg.AccessKey,
				SecretKey:   cfg.SecretKey,
				OpenKMS:     true,
				LogLevel:    string(cfg.LogLevel),
			},
		})
		if err != nil {
			return err
		}
		return nil
	}
	if err = retry.Retry(initNacos, 3, time.Second); err != nil {
		return nil, err
	}
	return
}
