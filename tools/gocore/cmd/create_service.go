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
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/template"
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

var dirList = []string{
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

var fileBuffer = new(bytes.Buffer)

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
	createConf(root, name)
	createMain(root, name)
	createDockerfile(root)
	createReadme(root)
	createErrCode(root)
	createModel(root, name)
	createCronjob(name, root)
	createJob(name, root)
	createApi(root, name)
	createDef(root)

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

func createMain(root, name string) {
	var cmdList []string
	if configJson.Get("api").Exists() {
		cmdList = append(cmdList, "cmd.Api,")
	}
	if configJson.Get("cronjob").Exists() {
		cmdList = append(cmdList, "cmd.Cronjob,")
	}
	if configJson.Get("job").Exists() {
		cmdList = append(cmdList, "cmd.Job,")
	}
	template.FromMain(name, cmdList, fileBuffer)
	fileWriter(fileBuffer, root+"/main.go")
}

func createConf(root string, name string) {

	// TODO: 如果用户不适用nacos就不应该初始化
	template.FromConfBase(fileBuffer)
	fileWriter(fileBuffer, root+"/conf/base.go")

	template.FromConfNacos(fileBuffer)
	fileWriter(fileBuffer, root+"/conf/nacos.go")

	template.FromConfConst(name, fileBuffer)
	fileWriter(fileBuffer, root+"/conf/const.go")
}

func createDockerfile(root string) {
	template.FromDockerfile(fileBuffer)
	fileWriter(fileBuffer, root+"/Dockerfile")
}

func createReadme(root string) {
	template.FromREADME(fileBuffer)
	fileWriter(fileBuffer, root+"/README.md")
}

func createErrCode(root string) {
	template.FromErrCode(fileBuffer)
	fileWriter(fileBuffer, root+"/app/errcode/errcode.go")
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
			template.FromModelTable(k1, tableStruct, tableName, fieldStr, fileBuffer)
			fileWriter(fileBuffer, tabelPath)

		}

		template.FromModel(k1, tableStr, fileBuffer)
		fileWriter(fileBuffer, dir+"/mysql_client.go")

		template.FromConfMysql(k1, fileBuffer)
		localConf += fileBuffer.String()
		template.FromConfLocal(localConf, fileBuffer)
		fileWriter(fileBuffer, root+"/conf/local.go")

		template.FromCmdInit(name, pkgs, dbUpdate, initDb, fileBuffer)
		fileWriter(fileBuffer, root+"/cmd/init.go")
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
		template.FromCronJob(k1, fileBuffer)
		fileWriter(fileBuffer, jobPath)
		cronjobs += "_ = cronObj.AddFunc(\"" + v1.String() + "\", cronjob." + k1 + ")\n"
	}

	template.FromCmdCronJob(name, cronjobs, fileBuffer)
	fileWriter(fileBuffer, root+"/cmd/cronjob.go")
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
		template.FromCronJob(k1, fileBuffer)
		fileWriter(fileBuffer, dir+file.CamelToUnderline(k1)+".go")
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

	template.FromCmdJob(name, jobCmd, jobFunctions, fileBuffer)
	fileWriter(fileBuffer, root+"/cmd/job.go")
}

func createApi(root, name string) {

	apiMap := configJson.Get("api").Map()
	if len(apiMap) == 0 {
		return
	}
	handlersList := apiMap["handlers"].Array()
	if len(handlersList) == 0 {
		return
	}

	template.FromCmdApi(name, fileBuffer)
	fileWriter(fileBuffer, root+"/cmd/api.go")

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

	routesStr := ""
	pkg := ""
	handlers := make([]string, 0)

	for _, v1 := range handlersList {
		handlerName := v1.Get("name").String()
		routesStr += "\n" + handlerName + ":=e.Group(\"" + v1.Get("prefix").String() + "\")\n"
		apiPath := apiDir + file.CamelToUnderline(handlerName) + ".go"
		routes := v1.Get("routes").Array()
		if len(routes) == 0 {
			continue
		}
		// 首字母大写
		handler := strings.Title(handlerName)
		handlers = append(handlers, handler)
		functions := make([]string, 0)
		reqs := make([]string, 0)
		for _, v2 := range routes {
			routesObj := strings.Split(v2.String(), ";")
			if len(routesObj) < 3 {
				continue
			}
			pkg = "\"" + name + "/app/api\"\n"
			route := routesObj[0]
			req := strings.Title(routesObj[1])
			reqs = append(reqs, req)
			function := strings.Title(route)
			functions = append(functions, function)
			routesStr += handlerName + ".POST(\"/" + handlerName + "/" + route + "\",api." + handler + "Handler." + function + ")\n"
			template.FromDomain(name, handler, function, req, fileBuffer)
			fileWriter(fileBuffer, domainDir+file.CamelToUnderline(route)+".go")

		}

		template.FromApi(name, handler, functions, reqs, fileBuffer)
		writer.Add(fileBuffer.Bytes())
		fileWriter(fileBuffer, apiPath)
	}

	template.FromApiRoutes(pkg, routesStr, fileBuffer)
	fileWriter(fileBuffer, root+"/app/routes/routers.go")

	template.FromDomainHandler(handlers, fileBuffer)
	fileWriter(fileBuffer, domainDir+"handler.go")
}

func createDef(root string) {
	structs := configJson.Get("api.structs").Map()
	if len(structs) == 0 {
		return
	}
	dir := root + "/app/def"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}

	writer.Add([]byte(`package def` + "\n"))
	for k1, v1 := range structs {
		params := ""
		fields := v1.Array()
		for _, v2 := range fields {
			field := strings.Split(v2.String(), ";")
			if len(field) < 3 {
				continue
			}
			params += file.UnderlineToCamel(field[0]) + " " + field[1] + " `json:\"" + field[0] + "\"`\n"
		}
		template.FromApiRequest(k1, params, fileBuffer)
	}
	fileWriter(fileBuffer, dir+"/def.go")
}

// ------------------------------------------------------------------------------

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

func fileWriter(buffer *bytes.Buffer, path string) {
	writer.Add(buffer.Bytes())
	writer.ForceWriteToFile(path)
	buffer.Reset()
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
