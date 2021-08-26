package cmd

import (
	"log"
	"os/exec"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/template"
	"gopkg.in/yaml.v2"

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

func CreatYoml(dir string, config *conf.GoCore) error {
	yamlPath := "gocore.yaml"
	if dir != "" {
		yamlPath = dir + "/gocore.yaml"
	}

	var writer = file.NewWriter()
	yamlByte, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	writer.Add(yamlByte)
	writer.WriteToFile(yamlPath)
	return nil
}

func creatService(c *cli.Context) error {
	config := conf.GetGocoreConfig()
	root := "."

	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = root
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	CreatYoml(root, config)
	template.CreateCode(root, config.Service.ProjectName, config)
	cmd = exec.Command("go", "mod", "init", config.Service.ProjectName)
	cmd.Dir = root
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("go", "mod", "tidy")
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

	cmd = exec.Command("golangci-lint", "run", "--exclude-use-default ")
	cmd.Dir = root
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	log.Println(config.Service.ProjectName + " 已生成...")
	return nil
}
