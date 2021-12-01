package cmd

import (
	"io/ioutil"
	"os"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// CreatYaml 创建配置文件
var CreatYaml = &cli.Command{
	Name: "yaml",
	Subcommands: []*cli.Command{
		{
			Name: "create",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dir",
					Usage:       "dir path",
					DefaultText: ".",
				}},
			Usage:  "create conf [dir]",
			Action: creatYaml,
		},
	},
}

// creatYaml 创建配置文件
func creatYaml(c *cli.Context) error {
	root := c.String("dir")
	if root == "" {
		root = "."
	}
	yamlPath := root + "/gocore.yaml"
	_, err := InitYaml(yamlPath, conf.GetGocoreConfig())
	if err != nil {
		return err
	}
	printHint("Welcome to GoCore, Configuration file has been generated.")
	return nil
}

// InitYaml 初始化Yaml配置文件
func InitYaml(yamlPath string, config *conf.GoCore) (*conf.GoCore, error) {
	if file.CheckFileIsExist(yamlPath) {
		apiFile, err := os.Open(yamlPath)
		if err == nil {
			content, err := ioutil.ReadAll(apiFile)
			if err != nil {
				panic(err)
			}
			cfg := conf.GoCore{}
			err = yaml.Unmarshal(content, &cfg)
			if err != nil {
				panic(err)
			}
			return &cfg, nil
		}
		panic(err)
	}

	return CreateYaml(yamlPath, config)
}

// CreateYaml 创建Yaml文件
func CreateYaml(yamlPath string, config *conf.GoCore) (*conf.GoCore, error) {
	var writer = file.NewWriter()
	yamlByte, err := yaml.Marshal(config)
	if err != nil {
		return config, err
	}
	writer.Add(yamlByte)
	writer.ForceWriteToFile(yamlPath)
	return config, nil
}
