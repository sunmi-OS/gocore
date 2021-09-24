package cmd

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/template"

	"github.com/urfave/cli/v2"
)

// CreatService 创建服务
var CreatService = &cli.Command{
	Name: "service",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from yaml file",
			DefaultText: "",
		}},
	Usage:  "create update service [config]",
	Action: creatService,
}

// creatService 创建服务并创建初始化配置
func creatService(c *cli.Context) error {
	config := conf.GetGocoreConfig()
	yamlPath := c.String("config")
	root := "."
	if yamlPath == "" {
		yamlPath = root + "/gocore.yaml"
	} else if !file.CheckFileIsExist(yamlPath) {
		return fmt.Errorf("%s is not found", yamlPath)
	}
	// 创建配置&读取配置
	config, err := InitYaml(yamlPath, config)
	if err != nil {
		panic(err)
	}

	modPath := root + "/go.mod"
	if file.CheckFileIsExist(modPath) {
		cmd := exec.Command("go", "fmt", "./...")
		cmd.Dir = root
		resp, err := cmd.Output()
		if err != nil {
			fmt.Println(string(resp))
			panic(err)
		}
	} else {
		printHint("Run go mod init.")
		cmd := exec.Command("go", "mod", "init", config.Service.ProjectName)
		cmd.Dir = root
		resp, err := cmd.Output()
		if err != nil {
			fmt.Println(string(resp))
			panic(err)
		}
	}

	template.CreateCode(root, config.Service.ProjectName, config)

	printHint("Run go mod tidy.")
	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Wait()
	if err != nil {
		return err
	}
	cmd.Dir = root
	resp, err := cmd.Output()
	if err != nil {
		fmt.Println(string(resp))
		panic(err)
	}

	printHint("Run go fmt.")
	cmd = exec.Command("go", "fmt", "./...")
	cmd.Dir = root
	resp, err = cmd.Output()
	if err != nil {
		fmt.Println(string(resp))
		panic(err)
	}
	printHint("Welcome to GoCore, the project has been format.")

	printHint("goimports -l -w .")
	cmd = exec.Command("goimports", "-l", "-w", ".")
	cmd.Dir = root
	resp, err = cmd.Output()
	if err != nil {
		fmt.Println(string(resp))
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
