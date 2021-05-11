package template

import "strings"

const (
	OneBackQuote = "`"
	TwoBackQuote = "``"
)

func CreateCmdApi(name string) string {
	return `
package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"` + name + `/app/routes"
	"` + name + `/common"

	"github.com/sunmi-OS/gocore/aliyunlog"
	"github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/web"
	"github.com/urfave/cli"
)

var Api = cli.Command{
	Name:    "api",
	Aliases: []string{"a"},
	Usage:   "api start",
	Subcommands: []cli.Command{
		{
			Name:   "start",
			Usage:  "开启运行api服务",
			Action: RunApi,
		},
	},
}

func RunApi(c *cli.Context) error {
	initConf()
	common.Init()
	initDB()
	e := web.InitEcho(&web.Config{
		Port: viper.C.GetString("network.ApiServicePort"),
	})
	routes.Router(e.Echo)

	e.Start()
	// 监听信号
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Fatalf("get a signal %s, stop the process", si.String())
			// Close相关服务
			e.Echo.Close()
			gorm.Close()
			aliyunlog.Close()
			return nil
		case syscall.SIGHUP:
		default:
			return nil
		}
	}
}
`
}
func CreateDockerfile() string {
	return `
#template
FROM sunmi-docker-images-registry.cn-hangzhou.cr.aliyuncs.com/public/golang:1.15 As builder

ENV GOPROXY https://mirrors.aliyun.com/goproxy/,direct
ENV GO111MODULE on

#step 1 build go cache
WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

#step 2 build binary project
WORKDIR /project
ADD . .
RUN ls
RUN go build main.go

FROM sunmi-docker-images-registry.cn-hangzhou.cr.aliyuncs.com/public/centos:7.8.2003
#run binary project
WORKDIR /app
COPY --from=builder /project/main .

# your project shell [project] [arg1] [arg2] ...
CMD [ "/app/main","api","start"]
`
}

func CreateReadme() string {
	return `
## 项目名称
> 请介绍一下你的项目吧  



## 运行条件
> 列出运行该项目所必须的条件和相关依赖  
* 条件一
* 条件二
* 条件三



## 运行说明
> 说明如何运行和使用你的项目，建议给出具体的步骤说明
* 操作一
* 操作二
* 操作三  



## 测试说明
> 如果有测试相关内容需要说明，请填写在这里  



## 技术架构
> 使用的技术框架或系统架构图等相关说明，请填写在这里  


## 协作者
> 高效的协作会激发无尽的创造力，将他们的名字记录在这里吧
`
}

func CreateToml() string {
	return `[service]
name = "order"

[api]
[api.publicOrder]
prefix = "/public/v1/order/"
routes = [
    "createPreOrder",
    "getPreOrder"
]
[api.privateOrder]
prefix = "/private/v1/order/"
routes = [
    "createPrivatePreOrder",
    "getPrivatePreOrder"
]

[cronjob]
StatisticDataByDay = "30 1 0 * * *"
LoopCSync = "30 1 0 * * *"

[job]
LoopOrder = "loopOrder"
LoopInvoice = "loopOrder"

[mysql]
[mysql.order]
tables = ["order","goods"]
[mysql.wallet]
tables = ["record"]

[redis]
[redis.order]
`
}

func CreateMain(name, cmds string) string {
	return `
package main

import (
	"log"
	"os"

	"` + name + `/cmd"
	"` + name + `/common"

	gocoreLog "github.com/sunmi-OS/gocore/log"
	"github.com/urfave/cli"
)

func main() {
	// 配置cli参数
	app := cli.NewApp()
	app.Name = common.PROJECT_NAME
	app.Usage = common.PROJECT_NAME
	app.Email = ""
	app.Version = common.PROJECT_VERSION

	// 指定命令运行的函数
	app.Commands = []cli.Command{
` + cmds + `
	}

	//初始化log
	gocoreLog.InitLogger("` + name + `")

	// 启动cli
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
`
}

func CreateConfBase() string {
	return `
package conf

var baseConfig = ` + OneBackQuote + `
[network]
ApiServiceHost = ""
ApiServicePort = "80"
` + OneBackQuote
}

func CreateConfMyql(dbName string) string {
	return `
[db` + strings.Title(dbName) + `]
dbHost = ""           #数据库连接地址
dbName = "` + dbName + `"           #数据库名称
dbUser = ""           #数据库用户名
dbPasswd = ""         #数据库密码
dbPort = "3306"       #数据库端口号
dbOpenconns_max = 20  #最大连接数
dbIdleconns_max = 20  #最大空闲连接
dbType = "mysql"
		`
}

func CreateConfLocal(content string) string {
	return `
package conf

var localConfig = ` + OneBackQuote + content + OneBackQuote
}

func CreateConfNacos() string {
	return `
package conf

import (
	"os"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/sunmi-OS/gocore/nacos"
)

// InitNacos  通过环境变量初始化配置中心，从环境变量获取连接配置中心的必要配置
func InitNacos(runtime string) {
	nacos.SetRunTime(runtime)
	nacos.ViperTomlHarder.SetviperBase(baseConfig)
	switch runtime {
	case "local":
		nacos.AddLocalConfig(runtime, localConfig)
	default:
		Endpoint := os.Getenv("ENDPOINT")
		NamespaceId := os.Getenv("NAMESPACE_ID")
		AccessKey := os.Getenv("ACCESS_KEY")
		SecretKey := os.Getenv("SECRET_KEY")

		if Endpoint == "" || NamespaceId == "" || AccessKey == "" || SecretKey == "" {
			panic("The configuration file cannot be empty.")
		}

		err := nacos.AddAcmConfig(runtime, constant.ClientConfig{
			Endpoint:    Endpoint,
			NamespaceId: NamespaceId,
			AccessKey:   AccessKey,
			SecretKey:   SecretKey,
		})
		if err != nil {
			panic(err)
		}
	}
}
`
}

func CreateRoutes(pkg, routes string) string {
	return `
package routes

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	` + pkg + `
	"github.com/labstack/echo/v4"
	"github.com/sunmi-OS/gocore/log"
)

var (
	pid      int
	progname string
)

func init() {
	pid = os.Getpid()
	paths := strings.Split(os.Args[0], "/")
	paths = strings.Split(paths[len(paths)-1], string(os.PathSeparator))
	progname = paths[len(paths)-1]

}

// Router 初始化路由
func Router(e *echo.Echo) {

	// 内存溢出检测
	e.Any("/pprof-init", func(context echo.Context) error {
		pid = os.Getpid()
		paths := strings.Split(os.Args[0], "/")
		paths = strings.Split(paths[len(paths)-1], string(os.PathSeparator))
		progname = paths[len(paths)-1]
		runtime.MemProfileRate = 1
		return nil
	})
	// 内存溢出检测
	e.Any("/pprof", func(context echo.Context) error {
		runtime.GC()
		f, err := os.Create(fmt.Sprintf("./heap_%s_%d_%s.prof", progname, pid, time.Now().Format("2006_01_02_03_04_05")))
		if err != nil {
			return err
		}
		defer f.Close()
		err = pprof.Lookup("heap").WriteTo(f, 1)
		if err != nil {
			log.Sugar.Info(err)
		}
		runtime.MemProfileRate = 0
		return context.JSON(200, "pong")
	})

` + routes + `
}
`
}

func CreateApi(name, handler string, functions ...string) string {
	res := `
package api

import (
	"` + name + `/app/domain"
	"` + name + `/app/errcode"
	"` + name + `/pkg/parse"

	"github.com/labstack/echo/v4"
)

var ` + handler + `Handler = ` + handler + `{}
type ` + handler + ` struct{}
`
	for _, v1 := range functions {
		res += `
// ` + v1 + `
func (*` + handler + `) ` + v1 + `(c echo.Context) error {
	params := new(domain.` + v1 + `Request)
	//参数验证绑定
	_, response, err := parse.ParseJson(c, params)
	if err != nil {
		return response.RetError(err, errcode.Code0002)
	}
	resp, code, err := domain.` + handler + `Handler.` + v1 + `(params)
	if err != nil {
		return response.RetError(err, code)
	}
	return response.RetSuccess(resp)
}
`
	}

	return res
}

func CreateDomain(handler, function string) string {
	return `
package domain

// ` + function + `
func (this *` + handler + `) ` + function + `(req *` + function + `Request)(map[string]interface{}, int, error) {
	return map[string]interface{}{}, 1, nil
}
`
}

func CreateDomainHandler(handlers ...string) string {
	res := `
package domain
`
	for _, v1 := range handlers {
		res += `var ` + v1 + `Handler = &` + v1 + `{}
		type ` + v1 + ` struct{}
`
	}
	return res

}

func CreateDomainRequest(requests ...string) string {
	res := `
package domain
`
	for _, v1 := range requests {
		res += "type " + v1 + `Request struct {
		OrderNo string ` + OneBackQuote + `json:"order_no"` + OneBackQuote + `
	}
	`
	}
	return res
}

func CreateCronjob(cron string) string {
	return `
	package cronjob
// ` + cron + `
func ` + cron + `() {
}
`
}

func CreateJob(job string) string {
	return `
	package job
// ` + job + `
func ` + job + `() {
}
`
}

func CreateCmdCronjob(name, cronjobs string) string {
	return `
package cmd

import (
	"log"
	"` + name + `/app/cronjob"
	"` + name + `/common"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
	"github.com/sunmi-OS/gocore/aliyunlog"
	"github.com/sunmi-OS/gocore/gorm"
	"github.com/urfave/cli"
)

// Cronjob cmd 定时任务相关
var Cronjob = cli.Command{
	Name:    "cron",
	Aliases: []string{"c"},
	Usage:   "run",
	Subcommands: []cli.Command{
		{
			Name:   "start",
			Usage:  "开启运行api服务",
			Action: runCron,
		},
	},
}

// runCron 运行定时任务
func runCron(c *cli.Context) error {

	// 初始化必要内容
	initConf()
	initDB()
	common.Init()
	cronObj := cron.New()
	
` + cronjobs + `

	// 同步阻塞运行
	cronObj.Start()

	// 监听信号
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Fatalf("get a signal %s, stop the process", si.String())
			// Close相关服务
			cronObj.Stop()
			gorm.Close()
			aliyunlog.Close()
			return nil
		case syscall.SIGHUP:
		default:
			return nil
		}
	}
}
	`
}

func CreateCmdJob(name, jobCmd, jobFunctions string) string {
	return `
package cmd

import (
	"` + name + `/app/job"
	"` + name + `/common"

	"github.com/urfave/cli"
)

// Job cmd 任务相关
var Job = cli.Command{
	Name:    "job",
	Aliases: []string{"j"},
	Usage:   "job",
	Subcommands: []cli.Command{
		` + jobCmd + `
	},
}
` + jobFunctions
}

func CreateModelClient(dbName string, tabels string) string {
	return `
package ` + dbName + `

import (
	"fmt"

	"github.com/jinzhu/gorm"
	g "github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/utils"
)

func Orm() *gorm.DB {
	db := g.GetORM("` + dbName + `")
	if utils.GetRunTime() != "onl" {
		db = db.Debug()
	}
	return db
}

func CreateTable() {
	fmt.Println("开始初始化` + dbName + `数据库")
	//自动建表，数据迁移
` + tabels + `
	fmt.Println("数据库` + dbName + `初始化完成")
}
`

}

func CreateModelTable(dbName, tableStruct, tableName string) string {
	return `
package ` + dbName + `

import (
	gormx "github.com/jinzhu/gorm"
)

var ` + tableStruct + `Handler = &` + tableStruct + `{}

type ` + tableStruct + ` struct {
}

func (* ` + tableStruct + `) TableName() string {
	return "` + tableName + `"
}

func (* ` + tableStruct + `) Insert(db *gormx.DB, data * ` + tableStruct + `) error {
	if db == nil {
		db = Orm()
	}
	return db.Create(data).Error
}

func (*` + tableStruct + `) GetOne(where string, args ...interface{}) (*` + tableStruct + `, error) {
	var obj ` + tableStruct + `
	return &obj, Orm().Where(where, args...).Take(&obj).Error
}

func (*` + tableStruct + `) GetList(where string, args ...interface{}) ([]*` + tableStruct + `, error) {
	var list []*` + tableStruct + `
	db := Orm()
	return list, db.Where(where, args...).Find(&list).Error
}

func (*` + tableStruct + `) GetCount(where string, args ...interface{}) (int, error) {
	var number int
	err := Orm().Model(&` + tableStruct + `{}).Where(where, args...).Count(&number).Error
	return number, err
}

func (*` + tableStruct + `) Delete(db *gormx.DB, where string, args ...interface{}) error {
	if db == nil {
		db = Orm()
	}
	return db.Where(where, args...).Delete(&` + tableStruct + `{}).Error
}

func (*` + tableStruct + `) Update(db *gormx.DB, data map[string]interface{}, where string, args ...interface{}) (int64, error) {
	if db == nil {
		db = Orm()
	}
	db = db.Model(&` + tableStruct + `{}).Where(where, args...).Update(data)
	return db.RowsAffected, db.Error
}
`
}

func CreateCommonConst(name string) string {
	return `package common
const (
	PROJECT_NAME    = "` + name + `"
	PROJECT_VERSION = "v1.0.0"
)
`
}

func CreateCmdInit(name, pkgs, dbUpdate, initDb string) string {
	return `
	package cmd

import (
	"log"

	` + pkgs + `
	"` + name + `/conf"

	"github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/nacos"
	"github.com/sunmi-OS/gocore/utils"
)

// initConf 初始化配置服务 （内部方法）
func initConf() {
	// 初始化Nacos配置
	conf.InitNacos(utils.GetRunTime())
	// 注册需要的配置
	nacos.ViperTomlHarder.SetDataIds("` + name + `", "mysql", "config", "redis")
	// 注册配置更新回调
	nacos.ViperTomlHarder.SetCallBackFunc("` + name + `", "mysql", func(namespace, group, dataId, data string) {
		` + dbUpdate + `
	})
	// 把Nacos的配置注册到Viper
	nacos.ViperTomlHarder.NacosToViper()
}

// initDB 初始化DB服务 （内部方法）
func initDB() {
	` + initDb + `
}
`
}

func CreateErrCode() string {
	return `
package errcode

const (
	Code0001 = iota + 11030001 //系统异常
	Code0002                   //参数错误
)

const CodeSuccess int64 = 1 //返回成功
`
}

func CreateParse() string {
	return `
package parse

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sunmi-OS/gocore/api"
)

func ParseJson(c echo.Context, req interface{}) (*api.Request, *api.Response, error) {
	request := api.NewRequest(c)
	response := api.NewResponse(c)
	err := request.InitRawJson()
	if err != nil {
		return request, response, err
	}
	request.GetRoot().GetJsonObject(req) //校验必填参数
	err = request.GetError()
	if err != nil {
		return request, response, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	return request, response, err
}

	`
}

func CreateCommon() string {
	return `package common
func Init(){
	
}
`
}
