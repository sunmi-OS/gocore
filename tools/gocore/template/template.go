package template

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/def"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
	"github.com/tidwall/gjson"
)

var writer = file.NewWriter()

var fileBuffer = new(bytes.Buffer)

var localConf string

var configJson gjson.Result

// hero -source=./tools/gocore/template -extensions=.got,.md,.docker

func CreateCode(root, name string, j gjson.Result) {
	configJson = j
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
	if configJson.Get("api").Exists() {
		cmdList = append(cmdList, "cmd.Api,")
	}
	if configJson.Get("cronjob").Exists() {
		cmdList = append(cmdList, "cmd.Cronjob,")
	}
	if configJson.Get("job").Exists() {
		cmdList = append(cmdList, "cmd.Job,")
	}
	FromMain(name, cmdList, fileBuffer)
	fileWriter(fileBuffer, root+"/main.go")
}

func createConf(root string, name string) {

	// TODO: 如果用户不适用nacos就不应该初始化
	FromConfBase(fileBuffer)
	fileWriter(fileBuffer, root+"/conf/base.go")

	FromConfNacos(fileBuffer)
	fileWriter(fileBuffer, root+"/conf/nacos.go")

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
				fieldStr += CreateField(v3.String())
			}
			FromModelTable(k1, tableStruct, tableName, fieldStr, fileBuffer)
			fileWriter(fileBuffer, tabelPath)

		}

		FromModel(k1, tableStr, fileBuffer)
		fileWriter(fileBuffer, dir+"/mysql_client.go")

		FromConfMysql(k1, fileBuffer)
		localConf += fileBuffer.String()
		FromConfLocal(localConf, fileBuffer)
		fileWriter(fileBuffer, root+"/conf/local.go")

		FromCmdInit(name, pkgs, dbUpdate, initDb, fileBuffer)
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
		FromCronJob(k1, fileBuffer)
		fileWriter(fileBuffer, jobPath)
		cronjobs += "_ = cronObj.AddFunc(\"" + v1.String() + "\", cronjob." + k1 + ")\n"
	}

	FromCmdCronJob(name, cronjobs, fileBuffer)
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
		FromCronJob(k1, fileBuffer)
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
	job.` + k1 + `()
	return nil
}
`
	}

	FromCmdJob(name, jobCmd, jobFunctions, fileBuffer)
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
			FromDomain(name, handler, function, req, fileBuffer)
			fileWriter(fileBuffer, domainDir+file.CamelToUnderline(route)+".go")

		}

		FromApi(name, handler, functions, reqs, fileBuffer)
		writer.Add(fileBuffer.Bytes())
		fileWriter(fileBuffer, apiPath)
	}

	FromApiRoutes(pkg, routesStr, fileBuffer)
	fileWriter(fileBuffer, root+"/app/routes/routers.go")

	FromDomainHandler(handlers, fileBuffer)
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
		FromApiRequest(k1, params, fileBuffer)
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
