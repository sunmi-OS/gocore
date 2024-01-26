package template

import "bytes"

func FromCmdApi(projectName string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package cmd

import (
	_ "`)
	buffer.WriteString(projectName)
	buffer.WriteString(`/errcode"
	"`)
	buffer.WriteString(projectName)
	buffer.WriteString(`/route"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/sunmi-OS/gocore/v2/api"

	"github.com/urfave/cli/v2"
)

var Api = &cli.Command{
	Name:    "api",
	Aliases: []string{"a"},
	Usage:   "api start",
	Subcommands: []*cli.Command{
		{
			Name:   "start",
			Usage:  "开启运行api服务",
			Action: RunApi,
		},
	},
}

func RunApi(c *cli.Context) error {
    // 初始化配置
    initConf()
	initDB()
	initCache()
	initLog()

    isDebugMode := true
	if utils.IsRelease() {
		isDebugMode = false
	}

	gs := api.NewGinServer(
	    api.WithServerDebug(isDebugMode),
	    api.WithServerHost(viper.C.GetString("network.ApiServiceHost")),
	    api.WithServerPort(viper.C.GetInt("network.ApiServicePort")),
	    api.WithOpenTrace(true),
    )
    // init route
    route.Routes(gs.Gin)
    gs.Start()

	return nil
}`)

}
