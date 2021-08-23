package template

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/def"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
	"github.com/tidwall/gjson"
)

var writer = file.NewWriter()

var fileBuffer = new(bytes.Buffer)

var localConf string

var configJson gjson.Result

var goCoreConfig *conf.GoCore

// hero -source=./tools/gocore/template -extensions=.got,.md,.docker

func CreateCode(root, name string, config *conf.GoCore) {
	goCoreConfig = config
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
}

func CreateToml() string {
	return `
[service]
name = "gen"

[api]
[api.structs]
GetPreOrderRequest = [
    "name;string;用户姓名",
    "dId;int64;用户dId",
]
CreatePreOrderRequest = [
    "name;string;用户姓名",
    "dId;int64;用户dId",
    "create_pre_order_content;struct:CreatePreOrderContent:详情"
]
CreatePreOrderContent = [
    "get_pre_order_content;*GetPreOrderContent;用户姓名",
    "list;[]*GetPreOrderContent;用户dId",
]
GetPreOrderContent = [
    "name;string;用户姓名",
    "dId;int64;用户dId",
]

[[api.handlers]]
name = "PublicOrder"
prefix = "/public/v1/order"
routes = [
    "createPreOrder;CreatePreOrderRequest;创建订单",
    "getPreOrder;GetPreOrderRequest;获取订单详情",
]

[[api.handlers]]
name = "PrivateOrder"
prefix = "/private/v1/order"
routes = [
    "createPrivatePreOrder;CreatePreOrderRequest;创建私有订单",
    "getPrivatePreOrder;GetPreOrderRequest;获取私有订单"
    ]
 
[cronjob]
StatisticDataByDay = "30 1 0 * * *"
LoopCSync = "30 1 0 * * *"

[job]
LoopOrder = "loopOrder"
LoopInvoice = "loopOrder"

[mysql]
[mysql.order]
order = [
    "column:id;primary_key;type:int AUTO_INCREMENT",
    "column:order_no;type:varchar(100) NOT NULL;default:'';comment:'订单号';unique_index",
    "column:uId;type:int NOT NULL;default:0;comment:'用户ID号';index",
    ]
goods = [
    "column:id;primary_key;type:int AUTO_INCREMENT",
    "column:order_no;type:varchar(100) NOT NULL;default:'';comment:'订单号';unique_index",
    "column:uId;type:int NOT NULL;default:0;comment:'用户ID号';index",
    "column:goods_id;type:varchar(50) NOT NULL;default:'';comment:'商品id';index",
    ]
[mysql.wallet]
record = [
    "column:id;primary_key;type:int AUTO_INCREMENT",
    "column:order_no;type:varchar(100) NOT NULL;default:'';comment:'订单号';unique_index",
    "column:uId;type:int NOT NULL;default:0;comment:'用户ID号';index",
    "column:goods_id;type:varchar(50) NOT NULL;default:'';comment:'商品id';index",
    "column:goods_num;type:int NOT NULL;default:'0';comment:'数量(sku属性)'",
    ]

[redis]
[redis.order]
`
}

func CreateField(field string) string {
	tags := strings.Split(field, ";")
	if len(tags) == 0 {
		return ""
	}

	fieldMap := make(map[string]string)
	for _, v1 := range tags {
		attributes := strings.Split(v1, ":")
		if len(attributes) < 2 {
			continue
		}
		fieldMap[attributes[0]] = attributes[1]
	}
	fieldName := fieldMap["column"]
	upFieldName := file.UnderlineToCamel(fieldName)
	fieldType := def.GetTypeName(fieldMap["type"])
	return upFieldName + "  " + fieldType + " `json:\"" + fieldName + "\" gorm:\"" + field + "\"`\n"
}

func createMain(root, name string) {
	var cmdList []string
	if goCoreConfig.HttpApiEnable {
		cmdList = append(cmdList, "cmd.Api,")
	}
	if goCoreConfig.CronJobEnable {
		cmdList = append(cmdList, "cmd.Cron,")
	}
	if goCoreConfig.JobEnable {
		cmdList = append(cmdList, "cmd.Job,")
	}
	FromMain(name, cmdList, fileBuffer)
	fileWriter(fileBuffer, root+"/main.go")
}

func createConf(root string, name string) {

	// TODO: 如果用户不适用nacos就不应该初始化
	FromConfBase(fileBuffer)
	fileWriter(fileBuffer, root+"/conf/base.go")

	// FromConfNacos(fileBuffer)
	// fileWriter(fileBuffer, root+"/conf/nacos.go")

	FromConfConst(name, fileBuffer)
	fileWriter(fileBuffer, root+"/conf/const.go")
}

func createDockerfile(root string) {
	FromDockerfile(fileBuffer)
	fileWriter(fileBuffer, root+"/Dockerfile")
}

func createReadme(root string) {
	FromREADME(fileBuffer)
	fileWriter(fileBuffer, root+"/README.md")
}

func createErrCode(root string) {
	FromErrCode(fileBuffer)
	fileWriter(fileBuffer, root+"/app/errcode/errcode.go")
}

func createModel(root, name string) {
	mysqlMap := goCoreConfig.Config.CMysql
	if len(mysqlMap) == 0 {
		return
	}
	pkgs := ""
	dbUpdate := ""
	if len(mysqlMap) > 0 {
		dbUpdate = "var err error"
	}
	initDb := ""
	for _, v1 := range mysqlMap {
		pkgs += `"` + name + `/app/model/` + v1.Name + `"` + "\n"
		dir := root + "/app/model/" + v1.Name
		dbUpdate += `
				err = orm.NewOrUpdateDB(conf.DB` + strings.Title(v1.Name) + `)
				if err != nil {
					glog.Error(err)
				}
		`
		initDb += `orm.NewDB(conf.DB` + strings.Title(v1.Name) + `)
			` + v1.Name + `.SchemaMigrate()` + "\n"
		err := file.MkdirIfNotExist(dir)
		if err != nil {
			panic(err)
		}
		tables := v1.Models
		tableStr := ""

		for _, v2 := range tables {
			tableName := v2.Name
			tableStruct := file.UnderlineToCamel(v2.Name)
			tableStr += "Orm().Set(\"gorm:table_options\", \"CHARSET=utf8mb4 comment='' AUTO_INCREMENT=1;\").AutoMigrate(&" + tableStruct + "{})\n"
			tabelPath := dir + "/" + tableName + ".go"
			fieldStr := ""
			fields := v2.Fields
			for _, v3 := range fields {
				fieldStr += CreateField(v3.GormRule)
			}
			FromModelTable(v1.Name, tableStruct, tableName, fieldStr, fileBuffer)
			fileWriter(fileBuffer, tabelPath)

		}

		FromModel(v1.Name, tableStr, fileBuffer)
		fileWriter(fileBuffer, dir+"/mysql_client.go")

		buff := new(bytes.Buffer)
		FromConfMysql(v1.Name, buff)
		localConf += buff.String()

		initRedis := ""
		for _, v1 := range goCoreConfig.Config.CRedis {
			for k2 := range v1.Index {
				localConf += `
[` + v1.Name + `]
host = "" 
port = ":6379"
auth = ""
prefix = ""
encryption = 0

[` + v1.Name + `.redisDB]
` + k2 + ` = ` + cast.ToString(v1.Index[k2])
				initRedis += "redis.NewRedis(conf." + strings.Title(v1.Name) + strings.Title(k2) + "Redis)\n"
				dbUpdate += `		
				err = redis.NewOrUpdateRedis(conf.` + strings.Title(v1.Name) + strings.Title(k2) + `Redis)
				if err != nil {
					glog.Error(err)
				}
		`
			}
		}
		if goCoreConfig.Config.CNacos.RocketMQConfig == true {
			localConf += `
[aliyunmq]
NameServer = ""
AccessKey = ""
SecretKey = ""
Namespace = ""

			`
		}
		FromConfLocal(localConf, fileBuffer)
		fileWriter(fileBuffer, root+"/conf/local.go")
		FromCmdInit(name, pkgs, dbUpdate, initDb+initRedis, fileBuffer)
		fileWriter(fileBuffer, root+"/cmd/init.go")
	}
}

func createCronjob(name, root string) {
	jobs := goCoreConfig.CronJobs
	if !goCoreConfig.CronJobEnable {
		return
	}

	dir := root + "/app/cronjob/"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}
	cronjobs := ""
	for _, v1 := range jobs {
		jobPath := dir + file.CamelToUnderline(v1.Job.Name) + ".go"
		FromCronJob(v1.Job.Name, fileBuffer)
		fileWriter(fileBuffer, jobPath)
		cronjobs += "_,_ = cronJob.AddFunc(\"" + v1.Spec + "\", cronjob." + v1.Job.Name + ")\n"
	}

	FromCmdCronJob(name, cronjobs, fileBuffer)
	fileWriter(fileBuffer, root+"/cmd/cron.go")
}

func createJob(name, root string) {

	jobs := goCoreConfig.Jobs
	if len(jobs) == 0 || !goCoreConfig.JobEnable {
		return
	}

	dir := root + "/app/job/"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}
	jobCmd := ""
	jobFunctions := ""
	for _, v1 := range jobs {
		FromJob(v1.Name, fileBuffer)
		fileWriter(fileBuffer, dir+file.CamelToUnderline(v1.Name)+".go")
		jobCmd += `		{
			Name:   "` + v1.Name + `",
			Usage:  "` + v1.Usage + `",
			Action: ` + v1.Name + `,
		},`
		jobFunctions += `
func ` + v1.Name + `(c *cli.Context) error {
	defer closes.Close()
	// 初始化必要内容
	initConf()
	initDB()
	job.` + v1.Name + `()
	return nil
}
`
	}

	FromCmdJob(name, jobCmd, jobFunctions, fileBuffer)
	fileWriter(fileBuffer, root+"/cmd/job.go")
}

func createApi(root, name string) {
	if !goCoreConfig.HttpApiEnable {
		return
	}
	handlersList := goCoreConfig.HttpApis.Apis
	if len(handlersList) == 0 {
		return
	}

	FromCmdApi(name, fileBuffer)
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

	handlers := make([]string, 0)

	for _, v1 := range handlersList {
		handlerName := v1.ModuleName
		routesStr += "\n" + handlerName + ":=router.Group(\"" + v1.Prefix + "\")\n"
		apiPath := apiDir + file.CamelToUnderline(handlerName) + ".go"
		routes := v1.Handle
		if len(routes) == 0 {
			continue
		}
		// 首字母大写
		handler := strings.Title(handlerName)
		handlers = append(handlers, handler)
		functions := make([]string, 0)
		reqs := make([]string, 0)
		for _, v2 := range routes {

			route := v2.Name
			req := strings.Title(v2.Name)
			reqs = append(reqs, req)
			function := strings.Title(route)
			functions = append(functions, function)
			routesStr += handlerName + "." + v2.Method + "(\"/" + file.CamelToUnderline(route) + "\",api." + function + ")\n"
			FromDomain(name, handler, function, req, fileBuffer)
			fileWriter(fileBuffer, domainDir+file.CamelToUnderline(route)+".go")

		}

		FromApi(name, handler, functions, reqs, fileBuffer)
		// writer.Add(fileBuffer.Bytes())
		fileWriter(fileBuffer, apiPath)
	}
	FromApiRoutes(name, routesStr, fileBuffer)
	fileWriter(fileBuffer, root+"/app/routes/routers.go")

}

func createDef(root string) {
	modules := goCoreConfig.HttpApis.Apis
	if len(modules) == 0 {
		return
	}
	dir := root + "/app/def"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}

	writer.Add([]byte(`package def` + "\n"))
	for _, v1 := range modules {
		for _, v2 := range v1.Handle {

			params := ""
			fields := v2.RequestParams
			for _, v3 := range fields {
				// field := strings.Split(v2.String(), ";")
				// if len(field) < 3 {
				// 	continue
				// }
				params += file.UnderlineToCamel(v3.Name) + " " + v3.Type + " `json:\"" + v3.Name + "\" validate:\"" + v3.Validate + "\"`\n"
			}
			FromApiRequest(v2.Name+"Request", params, fileBuffer)

			params = ""
			fields = v2.ResponseParams
			for _, v3 := range fields {
				params += file.UnderlineToCamel(v3.Name) + " " + v3.Type + " `json:\"" + v3.Name + "\" " + v3.Validate + "`\n"
			}
			FromApiRequest(v2.Name+"Response", params, fileBuffer)
		}
	}
	fileWriter(fileBuffer, dir+"/def.go")
}

// ------------------------------------------------------------------------------

func ParseToml(path string) {
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
	for _, dir := range dirList {
		dir = root + dir
		err := file.MkdirIfNotExist(dir)
		if err != nil {
			panic(err)
		}
	}
}
