package main

import (
	"log"
	"os"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/cmd"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	// 打印banner
	utils.PrintBanner(conf.PROJECT_NAME)
	// 配置cli参数
	app := cli.NewApp()
	app.Name = conf.PROJECT_NAME
	app.Usage = conf.PROJECT_NAME
	app.Version = conf.PROJECT_VERSION
	//app.Action = cmd.Ui.Action
	// 指定命令运行的函数
	app.Commands = []*cli.Command{
		cmd.CreatService,
		cmd.CreatYaml,
		cmd.AddMysql,
		cmd.Ui,
	}
	// 启动cli
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
