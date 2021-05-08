package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/tools/gocore/file"
	"github.com/sunmi-OS/gocore/tools/gocore/template"
	"github.com/tidwall/gjson"

	"github.com/urfave/cli"
)

// 创建服务
var CreatService = cli.Command{
	Name:  "create",
	Usage: "create cmd",
	Subcommands: []cli.Command{
		{
			Name: "service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Load configuration from toml file",
				}},
			Usage:  "create service [config]",
			Action: creatService,
		},
		{
			Name: "toml",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir",
					Usage: "dir path",
				}},
			Usage:  "create toml [dir]",
			Action: creatToml,
		},
	},
}

var dirList []string = []string{
	"/common",
	"/cmd",
	"/app/domain",
	"/app/model",
	"/app/apicode",
	"/app/routes",
	"/conf",
	"/pkg",
}

var writer = file.NewWriter()

var configJson gjson.Result

func creatService(c *cli.Context) error {
	config := c.String("config")
	if config == "" {
		return cli.NewExitError("config not found", 86)
	}
	parseToml(config)
	name := configJson.Get("service.name").String()
	if name == "" {
		return cli.NewExitError("service name  not found", 86)
	}
	name = "."
	mkdir(name)
	createMain(name)
	createConf(name)
	createRoutes(name)
	createCmd(name)
	createCommon(name)
	createDockerfile(name)
	createReadme(name)
	createModel(name)
	createJob(name)
	createApi(name)
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = name
	cmd.Output()

	cmd = exec.Command("go", "test", "./...")
	cmd.Dir = name
	cmd.Output()
	exec.Command("goimports", "-l", "-w", name).Output()
	log.Println(name + " 已生成...")
	return nil
}

func createMain(root string) {
	mainPath := root + "/main.go"
	writer.Add([]byte(template.MainTemplate))
	writer.WriteToFile(mainPath)
}

func createConf(root string) {
	confBasePath := root + "/conf/base.go"
	writer.Add([]byte(template.ConfBaseTemplate))
	writer.WriteToFile(confBasePath)

	confLocalPath := root + "/conf/local.go"
	writer.Add([]byte(template.ConfLocalTemplate))
	writer.WriteToFile(confLocalPath)

	confNacosPath := root + "/conf/nacos.go"
	writer.Add([]byte(template.ConfNacosTemplate))
	writer.WriteToFile(confNacosPath)
}

func createRoutes(root string) {
	routesPath := root + "/app/routes/routers.go"
	writer.Add([]byte(template.RoutesTemplate))
	writer.WriteToFile(routesPath)
}

func createCmd(root string) {
	cmdInitPath := root + "/cmd/init.go"
	writer.Add([]byte(strings.Join(template.CmdInitTemplate, root)))
	writer.WriteToFile(cmdInitPath)
}

func createCommon(root string) {
	commonPath := root + "/common/common.go"
	writer.Add([]byte(template.CommonTemplate))
	writer.WriteToFile(commonPath)
}

func createDockerfile(root string) {
	dockerFilePath := root + "/Dockerfile"
	writer.Add([]byte(template.DockerfileTemplate))
	writer.WriteToFile(dockerFilePath)
}

func createReadme(root string) {
	readmePath := root + "/README.md"
	writer.Add([]byte(template.ReadmeTemplate))
	writer.WriteToFile(readmePath)
}

func createModel(root string) {
	mysqlMap := configJson.Get("mysql").Map()
	if len(mysqlMap) == 0 {
		return
	}
	for k1, v1 := range mysqlMap {
		dir := root + "/app/model/" + k1
		err := file.MkdirIfNotExist(dir)
		if err != nil {
			panic(err)
		}

		clientPath := dir + "/client.go"
		writer.Add([]byte(template.ModelClientTemplate))
		writer.WriteToFile(clientPath)
		tables := v1.Get("tables").Array()
		for _, v2 := range tables {
			tabelPath := dir + "/" + file.CamelToUnderline(cast.ToString(v2)) + ".go"
			writer.Add([]byte(template.ModelTableTemplate))
			writer.WriteToFile(tabelPath)
		}

	}
}

func createJob(root string) {
	jobs := configJson.Get("job").Map()
	if len(jobs) == 0 {
		return
	}
	jobCmdPath := root + "/cmd/job.go"
	writer.Add([]byte(template.JobCmdTemplate))
	writer.WriteToFile(jobCmdPath)

	dir := root + "/app/job/"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}

	for k1, _ := range jobs {
		jobPath := dir + file.CamelToUnderline(k1) + ".go"
		writer.Add([]byte(template.JobTemplate))
		writer.WriteToFile(jobPath)

	}
}

func createApi(root string) {
	apiMap := configJson.Get("api").Map()
	if len(apiMap) == 0 {
		return
	}
	cmdApiPath := root + "/cmd/api.go"
	writer.Add([]byte(template.CmdApiTemplate))
	writer.WriteToFile(cmdApiPath)

	apiDir := root + "/app/api/"
	err := file.MkdirIfNotExist(apiDir)
	if err != nil {
		panic(err)
	}

	domainDir := root + "/app/domain/"
	err = file.MkdirIfNotExist(domainDir)
	if err != nil {
		panic(err)
	}
	apiWriter := file.NewWriter()
	domainWriter := file.NewWriter()
	for k1, v1 := range apiMap {
		apiPath := apiDir + file.CamelToUnderline(k1) + ".go"
		tables := v1.Get("routes").Array()
		if len(tables) == 0 {
			continue
		}
		apiWriter.Add([]byte(template.ApiTemplate[0]))
		for _, v2 := range tables {
			apiWriter.Add([]byte(template.ApiTemplate[1]))
			domainPath := domainDir + file.CamelToUnderline(cast.ToString(v2)) + ".go"
			domainWriter.Add([]byte(template.DomainTemplate))
			domainWriter.WriteToFile(domainPath)

		}
		apiWriter.WriteToFile(apiPath)
	}
}

func mkdir(root string) {
	for _, dir := range dirList {
		dir = root + dir
		err := file.MkdirIfNotExist(dir)
		if err != nil {
			panic(err)
		}
	}
}

func parseToml(path string) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	tree, err := toml.LoadReader(bytes.NewBuffer(buf))
	if err != nil {
		panic(err)
	}
	cMap := tree.ToMap()
	cMapBytes, err := json.Marshal(cMap)
	if err != nil {
		panic(err)
	}
	configJson = gjson.ParseBytes(cMapBytes)
}
