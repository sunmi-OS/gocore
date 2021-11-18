package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/template"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/urfave/cli/v2"
)

// CreatService 创建服务
var CreatService = &cli.Command{
	Name: "service",
	Subcommands: []*cli.Command{
		{
			Name: "create",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "config, c",
					Usage:       "Load configuration from yaml file",
					DefaultText: "",
				}},
			Usage:  "create update service [config]",
			Action: creatService,
		},
	},
}

// creatService 创建服务并创建初始化配置
func creatService(c *cli.Context) error {
	config := conf.GetGocoreConfig()
	yamlPath := c.String("config")
	root := "."

	if yamlPath == "" {
		yamlPath = root + "/gocore.yaml"
	}

	if !file.CheckFileIsExist(yamlPath) {
		return fmt.Errorf("%s is not found", yamlPath)
	}

	// 创建配置&读取配置
	config, err := InitYaml(yamlPath, config)
	if err != nil {
		panic(err)
	}

	modPath := root + "/go.mod"
	if file.CheckFileIsExist(modPath) {
		resp, err := utils.Cmd("go", []string{"fmt", "./..."})
		if err != nil {
			fmt.Println(resp)
			panic(err)
		}
	} else {
		printHint("Run go mod init.")
		resp, err := utils.Cmd("go", []string{"mod", "init", config.Service.ProjectName})
		if err != nil {
			fmt.Println(resp)
			panic(err)
		}
	}

	template.CreateCode(root, config.Service.ProjectName, config)

	printHint("Run go mod tidy.")

	resp, err := utils.Cmd("go", []string{"mod", "tidy"})
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	printHint("Run go fmt.")
	resp, err = utils.Cmd("go", []string{"fmt", "./..."})
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	printHint("goimports -l -w .")
	resp, err = utils.Cmd("goimports", []string{"-l", "-w", "."})
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}
	printHint("Welcome to GoCore, the project has been initialized.")

	return nil
}

// printHint 打印提示
func printHint(str string) {
	_, err := color.New(color.FgCyan, color.Bold).Print(str + "\n")
	if err != nil {
		return
	}
}
