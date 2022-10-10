package cmd

import (
	"fmt"
	"os"
	"os/exec"

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
	sourceCodeRoot := root + "/app"
	err := os.Mkdir(sourceCodeRoot, os.ModePerm)
	if err != nil {
		panic("Failed to create sourceCodeRoot folder: " + err.Error())
	}
	if yamlPath == "" {
		yamlPath = root + "/gocore.yaml"
	}

	if !file.CheckFileIsExist(yamlPath) {
		return fmt.Errorf("%s is not found", yamlPath)
	}

	// 创建配置&读取配置
	config, err = InitYaml(yamlPath, config)
	if err != nil {
		panic(err)
	}

	modPath := sourceCodeRoot + "/go.mod"
	if file.CheckFileIsExist(modPath) {
		resp, err := utils.Cmd("go", []string{"fmt", "./..."})
		if err != nil {
			fmt.Println(resp)
			panic(err)
		}
	} else {
		printHint("Run go mod init")
		goModInitCmd := exec.Command("go", []string{"mod", "init", config.Service.ProjectName}...)
		goModInitCmd.Dir = sourceCodeRoot
		out, err := goModInitCmd.Output()
		if err != nil {
			fmt.Println(out)
			panic(err)
		}
	}

	template.CreateCode(root, sourceCodeRoot, config.Service.ProjectName, config)

	printHint("Run go mod tidy")

	goModTidyCmd := exec.Command("go", []string{"mod", "tidy"}...)
	goModTidyCmd.Dir = sourceCodeRoot
	goModTidyCmd.Stderr = os.Stderr
	err = goModTidyCmd.Start()
	if err != nil {
		panic("go mod tidy error: " + err.Error())
	}
	_ = goModTidyCmd.Wait()

	printHint("Run go fmt")

	goFmtCmd := exec.Command("go", []string{"fmt", "./..."}...)
	goFmtCmd.Dir = sourceCodeRoot
	goFmtCmd.Stderr = os.Stderr
	err = goFmtCmd.Start()
	if err != nil {
		panic("go fmt error: " + err.Error())
	}
	_ = goFmtCmd.Wait()

	printHint("Welcome to GoCore, the project has been initialized")

	return nil
}

// printHint 打印提示
func printHint(str string) {
	_, err := color.New(color.FgCyan, color.Bold).Print(str + "\n")
	if err != nil {
		return
	}
}
