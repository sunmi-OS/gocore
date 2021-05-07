package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/sunmi-OS/gocore/tools/gocore/constant"
	"github.com/sunmi-OS/gocore/tools/gocore/file"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/urfave/cli"
)

// 创建服务
var CreatService = cli.Command{
	Name:    "create",
	Aliases: []string{"create service"},
	Usage:   "create cmd",
	Subcommands: []cli.Command{
		{
			Name: "service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Load configuration from toml file",
				}},
			Usage:  "create service",
			Action: creatService,
		},
		{
			Name: "model",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Create model file",
				}},
			Usage:  "create model",
			Action: creatService,
		},
	},
}

func creatService(c *cli.Context) error {
	config := c.String("config")
	if config == "" {
		return cli.NewExitError("config not found", 86)
	}
	viper.NewConfig(config, config)
	service := viper.C.GetStringMap("service")
	name := viper.C.GetString("service.name")
	err := file.MkdirIfNotExist(name)
	if err != nil {
		panic(err)
	}
	path := name + "/main.go"
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	if e != nil {
		return e
	}
	defer f.Close()
	writer := file.NewWriter()
	writer.Add([]byte(constant.MainTemplate))
	err = writer.WriteToFile(f)
	if err != nil {
		panic(err)
	} else {
		exec.Command("goimports", "-l", "-w", path).Output()
		log.Println(path + " 已生成...")
	}
	fmt.Printf("%#v\n", service)
	return nil
}
