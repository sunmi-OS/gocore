package cmd

import (
	"log"
	"os/exec"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/template"

	"github.com/tidwall/gjson"

	"github.com/urfave/cli/v2"
)

// 创建服务
var CreatService = &cli.Command{
	Name:  "create",
	Usage: "create cmd",
	Subcommands: []*cli.Command{
		{
			Name: "service",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "config, c",
					Usage: "Load configuration from toml file",
				}},
			Usage:  "create service [config]",
			Action: creatService,
		},
		{
			Name: "toml",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "dir",
					Usage: "dir path",
				}},
			Usage:  "create toml [dir]",
			Action: creatToml,
		},
	},
}

var configJson gjson.Result

func creatService(c *cli.Context) error {
	config := c.String("config")
	if config == "" {
		return cli.NewExitError("config not found", 86)
	}
	template.ParseToml(config)
	name := configJson.Get("service.name").String()
	if name == "" {
		return cli.NewExitError("service name  not found", 86)
	}
	root := "."

	template.CreateCode(root, name, configJson)

	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = root
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("goimports", "-l", "-w", "./...")
	cmd.Dir = root
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("go", "test", "./...")
	cmd.Dir = root
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("go", "fmt", "./...")
	cmd.Dir = root
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	log.Println(name + " 已生成...")
	return nil
}
