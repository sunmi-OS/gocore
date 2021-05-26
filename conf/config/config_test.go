package config

import (
	"os"
	"testing"

	"github.com/sunmi-OS/gocore/utils/xlog"
)

const (
	Group  = "GroupName"
	DataId = "DataId"
)

var LocalConfig = `
[cfg]
name = "gocore"

[http]
host = ""
port = ":1234"
debug = true
`

var Conf = &ConfigStruct{}

type ConfigStruct struct {
	Cfg  *Cfg  `toml:"cfg"`
	Http *Http `toml:"http"`
}

type Cfg struct {
	Name string `toml:"name"`
}

type Http struct {
	Host  string `toml:"host"`
	Port  string `toml:"port"`
	Debug bool   `toml:"debug"`
}

func TestParseConfig(t *testing.T) {
	os.Setenv("RUN_TIME", "local")
	err := Parse(Group, DataId, LocalConfig, Conf)
	if err != nil {
		panic(err)
	}
	xlog.Debugf("Conf.Cfg:%+v", Conf.Cfg)
	xlog.Debugf("Conf.Http:%+v", Conf.Http)
}
