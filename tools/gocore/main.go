package main

import (
	"log"
	"os"

	cmd "github.com/sunmi-OS/gocore/tools/gocore/cmd"
	"github.com/sunmi-OS/gocore/tools/gocore/conf"
	"github.com/urfave/cli"
)

func main() {
	// 配置cli参数
	app := cli.NewApp()
	app.Name = conf.PROJECT_NAME
	app.Usage = conf.PROJECT_NAME
	app.Email = ""
	app.Version = conf.PROJECT_VERSION
	// 指定命令运行的函数
	app.Commands = []cli.Command{
		cmd.CreatService,
	}

	// 启动cli
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
