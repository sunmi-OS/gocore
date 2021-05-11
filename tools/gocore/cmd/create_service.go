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
	"/app/errcode",
	"/app/routes",
	"/conf",
	"/pkg",
}

var writer = file.NewWriter()

var configJson gjson.Result

var localConf string

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
	root := "."
	mkdir(root)
	createConf(root)
	createParse(root)
	createMain(root, name)
	createCmd(root)
	createCommon(root, name)
	createDockerfile(root)
	createReadme(root)
	createErrCode(root)
	createModel(root, name)
	createCronjob(name, root)
	createJob(name, root)
	createApi(root, name)

	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = root
	cmd.Output()

	cmd = exec.Command("goimports", "-l", "-w", "./...")
	cmd.Dir = root
	cmd.Output()

	cmd = exec.Command("go", "test", "./...")
	cmd.Dir = root
	cmd.Output()

	cmd = exec.Command("go", "fmt", "./...")
	cmd.Dir = root
	cmd.Output()

	log.Println(name + " 已生成...")
	return nil
}

func createMain(root, name string) {
	mainPath := root + "/main.go"
	cmds := ""
	if configJson.Get("api").Exists() {
		cmds += "cmd.Api,\n"
	}
	if configJson.Get("cronjob").Exists() {
		cmds += "cmd.Cronjob,\n"
	}
	if configJson.Get("job").Exists() {
		cmds += "cmd.Job,\n"
	}
	writer.Add([]byte(template.CreateMain(name, cmds)))
	writer.ForceWriteToFile(mainPath)
}

func createConf(root string) {
	confBasePath := root + "/conf/base.go"
	writer.Add([]byte(template.CreateConfBase()))
	writer.WriteToFile(confBasePath)

	confNacosPath := root + "/conf/nacos.go"
	writer.Add([]byte(template.CreateConfNacos()))
	writer.WriteToFile(confNacosPath)
}

func createCmd(root string) {
	// cmdInitPath := root + "/cmd/init.go"
	// writer.Add([]byte(strings.Join(template.CmdInitTemplate, root)))
	// writer.WriteToFile(cmdInitPath)
}

func createCommon(root, name string) {
	commonPath := root + "/common/common.go"
	writer.Add([]byte(template.CreateCommon()))
	writer.WriteToFile(commonPath)

	constPath := root + "/common/const.go"
	writer.Add([]byte(template.CreateCommonConst(name)))
	writer.WriteToFile(constPath)
}

func createDockerfile(root string) {
	dockerFilePath := root + "/Dockerfile"
	writer.Add([]byte(template.CreateDockerfile()))
	writer.WriteToFile(dockerFilePath)
}

func createReadme(root string) {
	readmePath := root + "/README.md"
	writer.Add([]byte(template.CreateReadme()))
	writer.WriteToFile(readmePath)
}

func createErrCode(root string) {
	errCodePath := root + "/app/errcode/errcode.go"
	writer.Add([]byte(template.CreateErrCode()))
	writer.WriteToFile(errCodePath)
}

func createModel(root, name string) {
	mysqlMap := configJson.Get("mysql").Map()
	if len(mysqlMap) == 0 {
		return
	}
	pkgs := ""
	dbUpdate := ""
	if len(mysqlMap) > 0 {
		dbUpdate = "var err error"
	}
	initDb := ""
	for k1, v1 := range mysqlMap {
		pkgs += `"` + name + `/app/model/` + k1 + `"` + "\n"
		dir := root + "/app/model/" + k1
		dbUpdate += `
				err = gorm.NewOrUpdateDB("db` + strings.Title(k1) + `")
				if err != nil {
					log.Fatalln(err)
				}
		`
		initDb += `gorm.NewDB("db` + strings.Title(k1) + `")
			` + k1 + `.CreateTable()` + "\n"
		err := file.MkdirIfNotExist(dir)
		if err != nil {
			panic(err)
		}
		tables := v1.Map()
		tableStr := ""

		for k2, v2 := range tables {
			tableName := cast.ToString(k2)
			tableStruct := file.UnderlineToCamel(tableName)
			tableStr += "Orm().Set(\"gorm:table_options\", \"CHARSET=utf8mb4 comment='中台订单记录表' AUTO_INCREMENT=1;\").AutoMigrate(&" + tableStruct + "{})\n"
			tabelPath := dir + "/" + tableName + ".go"
			fieldStr := ""
			fields := v2.Array()
			for _, v3 := range fields {
				fieldStr += template.CreateField(v3.String())
			}
			writer.Add([]byte(template.CreateModelTable(k1, tableStruct, tableName, fieldStr)))
			writer.ForceWriteToFile(tabelPath)

		}

		clientPath := dir + "/mysql_client.go"
		writer.Add([]byte(template.CreateModelClient(k1, tableStr)))
		writer.ForceWriteToFile(clientPath)

		localConf += template.CreateConfMyql(k1)
		confLocalPath := root + "/conf/local.go"
		writer.Add([]byte(template.CreateConfLocal(localConf)))
		writer.ForceWriteToFile(confLocalPath)

		cmdInitPath := root + "/cmd/init.go"
		writer.Add([]byte(template.CreateCmdInit(name, pkgs, dbUpdate, initDb)))
		writer.ForceWriteToFile(cmdInitPath)
	}
}

func createCronjob(name, root string) {
	jobs := configJson.Get("cronjob").Map()
	if len(jobs) == 0 {
		return
	}

	dir := root + "/app/cronjob/"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}
	cronjobs := ""
	for k1, v1 := range jobs {
		jobPath := dir + file.CamelToUnderline(k1) + ".go"
		writer.Add([]byte(template.CreateCronjob(k1)))
		writer.WriteToFile(jobPath)
		cronjobs += "_ = cronObj.AddFunc(\"" + v1.String() + "\", cronjob." + k1 + ")\n"
	}

	cronCmdPath := root + "/cmd/cronjob.go"
	writer.Add([]byte(template.CreateCmdCronjob(name, cronjobs)))
	writer.ForceWriteToFile(cronCmdPath)
}

func createJob(name, root string) {
	jobs := configJson.Get("job").Map()
	if len(jobs) == 0 {
		return
	}

	dir := root + "/app/job/"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}
	jobCmd := ""
	jobFunctions := ""
	for k1, v1 := range jobs {
		jobPath := dir + file.CamelToUnderline(k1) + ".go"
		writer.Add([]byte(template.CreateJob(k1)))
		writer.WriteToFile(jobPath)
		jobCmd += `		{
			Name:   "` + v1.String() + `",
			Usage:  "开启运行api服务",
			Action: ` + k1 + `,
		},`
		jobFunctions += `
func ` + k1 + `(c *cli.Context) error {
	// 初始化必要内容
	initConf()
	initDB()
	common.Init()
	job.` + k1 + `()
	return nil
}
`
	}

	jobCmdPath := root + "/cmd/job.go"
	writer.Add([]byte(template.CreateCmdJob(name, jobCmd, jobFunctions)))
	writer.ForceWriteToFile(jobCmdPath)
}

func createApi(root, name string) {
	apiMap := configJson.Get("api").Map()
	if len(apiMap) == 0 {
		return
	}
	cmdApiPath := root + "/cmd/api.go"
	writer.Add([]byte(template.CreateCmdApi(name)))
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

	routesStr := ""
	pkg := ""
	handlers := make([]string, 0)
	requests := make([]string, 0)
	for k1, v1 := range apiMap {
		routesStr += "\n" + k1 + ":=e.Group(\"" + v1.Get("prefix").String() + "\")\n"
		apiPath := apiDir + file.CamelToUnderline(k1) + ".go"
		routes := v1.Get("routes").Array()
		if len(routes) == 0 {
			continue
		}
		// 首字母大写
		handler := strings.Title(k1)
		handlers = append(handlers, handler)
		functions := make([]string, 0)

		for _, v2 := range routes {
			pkg = "\"" + name + "/app/api\"\n"
			route := cast.ToString(v2)
			function := strings.Title(route)
			functions = append(functions, function)
			requests = append(requests, function)
			routesStr += k1 + ".POST(\"/" + k1 + "/" + route + "\",api." + handler + "Handler." + function + ")\n"

			domainPath := domainDir + file.CamelToUnderline(route) + ".go"
			domainWriter.Add([]byte(template.CreateDomain(handler, function)))
			domainWriter.WriteToFile(domainPath)

		}
		apiWriter.Add([]byte(template.CreateApi(name, handler, functions...)))
		apiWriter.WriteToFile(apiPath)
	}
	domainRequestPath := domainDir + "request.go"
	writer.Add([]byte(template.CreateDomainRequest(requests...)))
	writer.ForceWriteToFile(domainRequestPath)

	routesPath := root + "/app/routes/routers.go"
	writer.Add([]byte(template.CreateRoutes(pkg, routesStr)))
	writer.ForceWriteToFile(routesPath)

	domainHandlerPath := domainDir + "handler.go"
	writer.Add([]byte(template.CreateDomainHandler(handlers...)))
	writer.ForceWriteToFile(domainHandlerPath)
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

func createParse(root string) {
	dir := root + "/pkg/parse"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}
	parsePath := dir + "/parse.go"
	writer.Add([]byte(template.CreateParse()))
	writer.WriteToFile(parsePath)
}
