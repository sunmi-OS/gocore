package template

const (
	BackQuote    = "``"
	MainTemplate = `
package main

import (

	gocoreLog "github.com/sunmi-OS/gocore/log"

)

func main() {

	//初始化log
	gocoreLog.InitLogger("order")
	
}
`
	ConfBaseTemplate = `
package conf

var baseConfig = ` + BackQuote

	ConfLocalTemplate = `
package conf

var localConfig = ` + BackQuote
	ConfNacosTemplate = `
package conf

import (
	"os"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/sunmi-OS/gocore/nacos"
)

// @title        初始化配置中心
// @description  通过环境初始化配置中心，从环境变量获取连接配置中心的必要配置
// @param        runtime  string  "当前运行的环境字符串 local/dev/test/uat/onl等"
// @return       无
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

	RoutesTemplate = `
package routes

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

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
}
`
)

var (
	CmdInitTemplate = []string{`
	package cmd

import (
	"log"

	"order/app/model"
	"order/app/model/misun"
	"order/conf"

	"github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/nacos"
	"github.com/sunmi-OS/gocore/utils"
)

// initConf 初始化配置服务 （内部方法）
func initConf() {
	// 初始化Nacos配置
	conf.InitNacos(utils.GetRunTime())
	// 注册需要的配置
	nacos.ViperTomlHarder.SetDataIds("`, `", "mysql", "config", "redis")
	nacos.ViperTomlHarder.SetDataIds("general", "general")
	// 注册配置更新回调
	nacos.ViperTomlHarder.SetCallBackFunc("`, `", "mysql", func(namespace, group, dataId, data string) {
		err := gorm.NewOrUpdateDB(conf.OrderDB)
		if err != nil {
			log.Fatalln(err)
		}
	})
	// 把Nacos的配置注册到Viper
	nacos.ViperTomlHarder.NacosToViper()
}

// initDB 初始化DB服务 （内部方法）
func initDB() {
}
`}
	CmdApiTemplate = `
package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"order/app/routes"
	"order/common"

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
			common.Monitor.Close(5)
			aliyunlog.Close()
			return nil
		case syscall.SIGHUP:
		default:
			return nil
		}
	}
}
`
	CommonTemplate = `package common`

	DockerfileTemplate = `
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
	ReadmeTemplate = `
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

	ModelClientTemplate = `
package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	g "github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/utils"
)

func Orm() *gorm.DB {
	db := g.GetORM("order")
	if utils.GetRunTime() != "onl" {
		db = db.Debug()
	}
	return db
}

func CreateTable() {
	fmt.Println("开始初始化order数据库")
	//自动建表，数据迁移
	Orm().Set("gorm:table_options", "CHARSET=utf8mb4 comment='订单表' AUTO_INCREMENT=1;").AutoMigrate(&Order{})
	fmt.Println("数据库order初始化完成")
}
`

	ModelTableTemplate = `
package model

import (
	gormx "github.com/jinzhu/gorm"
)

var AccountHandler = &Account{}

type Account struct {
}

func (*Account) TableName() string {
	return "wallet_account"
}

func (*Account) Insert(db *gormx.DB, data *Account) error {
	if db == nil {
		db = Orm()
	}
	return db.Create(data).Error
}

func (*Account) GetOne(where string, args ...interface{}) (*Account, error) {
	var obj Account
	return &obj, Orm().Where(where, args...).Take(&obj).Error
}

func (*Account) GetList(where string, args ...interface{}) ([]*Account, error) {
	var list []*Account
	db := Orm()
	return list, db.Where(where, args...).Find(&list).Error
}

func (*Account) GetCount(where string, args ...interface{}) (int, error) {
	var number int
	err := Orm().Model(&Account{}).Where(where, args...).Count(&number).Error
	return number, err
}

func (*Account) Delete(db *gormx.DB, where string, args ...interface{}) error {
	if db == nil {
		db = Orm()
	}
	return db.Where(where, args...).Delete(&Account{}).Error
}

func (*Account) Update(db *gormx.DB, data map[string]interface{}, where string, args ...interface{}) (int64, error) {
	if db == nil {
		db = Orm()
	}
	db = db.Model(&Account{}).Where(where, args...).Update(data)
	return db.RowsAffected, db.Error
}
`
	JobCmdTemplate = `
package cmd

import (
	"fmt"
	"log"
	orderCron "order/app/domain/cron"
	"order/app/model"
	"order/common"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
	"github.com/sunmi-OS/gocore/aliyunlog"
	"github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/utils"
	"github.com/urfave/cli"
)

// Cron cmd 定时任务相关
var Cron = cli.Command{
	Name:    "cron",
	Aliases: []string{"c"},
	Usage:   "run",
	Subcommands: []cli.Command{
		{
			Name:   "start",
			Usage:  "开启运行api服务",
			Action: RunCron,
		},
	},
}

//执行一次性订单统计任务执行
var OneTaskOrderStatistic = cli.Command{
	Name:  "orderstatistic",
	Usage: "按时间段统计入库",
	Subcommands: []cli.Command{
		{
			Name:   "exec",
			Usage:  "按时间段手动执行统计脚本",
			Action: OneTaskOrderStatisFun,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "startDate", //  ./main orderstatistic exec --startDate=2020-12-21 --endDate=2020-12-21 --businessKey=test
					Value: "",
					Usage: "开始日期格式2020-01-01",
				},
				cli.StringFlag{
					Name:  "endDate",
					Value: "",
					Usage: "结束日期格式2020-01-10",
				},
				cli.StringFlag{
					Name:  "businessKey",
					Value: "",
					Usage: "业务类型标识 可空指所有",
				},
			},
		},
	},
}

// RunCron 运行定时任务
func RunCron(c *cli.Context) error {

	// 初始化必要内容
	initConf()
	initDB()
	common.Init()
	job := cron.New()
	go orderCron.LoopMonitorHandler.LoopMonitorList()
	go orderCron.LoopMonitorHandler.LoopMonitorMsg()
	go orderCron.LoopCSyncHandler.LoopCSync(1) // c端同步业务组
	if utils.GetRunTime() == "onl" {
		go orderCron.LoopSyncNcHandler.Loop()
		go orderCron.LoopSyncNcHandler.LoopInvoice()
		go orderCron.LoopCSyncHandler.LoopCSync(2) // c端同步nc
	}
	go orderCron.LoopSyncBusinessHandler.Loop()
	_ = job.AddFunc("30 1 0 * * *", orderCron.LoopExpireHandler.Loop)
	_ = job.AddFunc("30 5 0 * * *", orderCron.LoopOrderStatisticHandler.StatisticDataByDay) // 0点5分30秒开始执行

	// 同步阻塞运行
	job.Start()

	// 监听信号
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Fatalf("get a signal %s, stop the process", si.String())
			// Close相关服务
			job.Stop()
			gorm.Close()
			common.Monitor.Close(5)
			aliyunlog.Close()
			return nil
		case syscall.SIGHUP:
		default:
			return nil
		}
	}
}

//OneTaskOrderStatisFun 手动执行统计订单数据脚本
func OneTaskOrderStatisFun(c *cli.Context) error {

	// 初始化必要内容
	initConf()
	initDB()
	common.Init()

	var startDate string = c.String("startDate")
	var endDate string = c.String("endDate")
	var businessKey string = c.String("businessKey")
	fmt.Println("startDate=", startDate, "endDate=", endDate, "businessKey=", businessKey)
	if startDate == "" || endDate == "" {
		fmt.Println("请输入开始和结束日期")
		return nil
	}
	if businessKey == "" {
		//查询所有的businessKey
		businessList, err := model.BusinessTypeHandler.GetList("")
		if err != nil {
			_ = aliyunlog.Error("order-StatisticDataByDay-err", map[string]string{
				"err": err.Error(),
				"msg": "model.BusinessTypeHandler.GetList err,脚本退出",
			})
		}
		for _, val := range businessList {
			_ = orderCron.LoopOrderStatisticHandler.OrderStatistic(val.BusinessKey, startDate, endDate)
		}
	} else {
		_ = orderCron.LoopOrderStatisticHandler.OrderStatistic(businessKey, startDate, endDate)
	}
	return nil
}
	`
	JobTemplate = `
	package cron

import (
	"time"

	"order/app/domain"
	"order/app/model"
	"order/common"

	"github.com/sunmi-OS/gocore/log"
	"github.com/weblazy/core/mapreduce"
)

var LoopCSyncHandler = &LoopCSync{}

type LoopCSync struct{}

// @desc  回调通知
// @auth liuguoqiang 2020-11-24
// @param
// @return
func (this *LoopCSync) LoopCSync(syncType int) {
	for {
		var indexId int64 = 0
		for {
			now := time.Now().Unix()
			tempTime := now - 86400 //24个小时内的订单,下一次回调时间大于当前时间,回调次数小于10
			list, err := model.COrderSyncHandler.GetList("create_time > ? and sync_status = 0  and next_sync_time < ? and sync_num < 10 and id > ? and sync_type = ?", tempTime, now, indexId, syncType)

			if err != nil {
				log.Sugar.Info(err)
				common.SendMonitorMsg(common.MONITOR_LIST_SYNC_NC, "\n订单通知NC:查询未通知订单失败:---\n"+err.Error())
				time.Sleep(time.Second * 10)
				continue
			}

			if len(list) == 0 {
				time.Sleep(time.Second * 10)
				break
			}
			//并发处理
			mapreduce.MapVoid(func(source chan<- interface{}) {
				for _, item := range list {
					indexId = item.Id
					source <- item
				}
			}, func(obj interface{}) {
				item := obj.(*model.COrderSync)
				domain.OrderHandler.CSync(item)
			})
			time.Sleep(time.Second * 10)
		}
	}
}
`

	ApiTemplate = []string{`
package api

import (
	"order/app/domain"
	"order/errcode"
	"order/pkg/parse"

	"github.com/labstack/echo/v4"
)

var InvoiceHandler = Invoice{}

type Invoice struct{}
`, `
//Notice 发票结果回传
func (*Invoice) Notice(c echo.Context) error {
	params := new(domain.NoticeRequest)
	//参数验证绑定
	_, response, err := parse.ParseJson(c, params)
	if err != nil {
		return response.RetError(err, errcode.Code0002)
	}
	resp, code, err := domain.InvoiceHandler.InvoiceNotice(params)
	if err != nil {
		return response.RetError(err, code)
	}
	return response.RetSuccess(resp)
}
`}

	DomainTemplate = `
package cron

import (
	"time"

	"order/app/domain"
	"order/app/model"
	"order/common"

	"github.com/sunmi-OS/gocore/log"
	"github.com/weblazy/core/mapreduce"
)

var LoopCSyncHandler = &LoopCSync{}

type LoopCSync struct{}

// @desc  回调通知
// @auth liuguoqiang 2020-11-24
// @param
// @return
func (this *LoopCSync) LoopCSync(syncType int) {
	for {
		var indexId int64 = 0
		for {
			now := time.Now().Unix()
			tempTime := now - 86400 //24个小时内的订单,下一次回调时间大于当前时间,回调次数小于10
			list, err := model.COrderSyncHandler.GetList("create_time > ? and sync_status = 0  and next_sync_time < ? and sync_num < 10 and id > ? and sync_type = ?", tempTime, now, indexId, syncType)

			if err != nil {
				log.Sugar.Info(err)
				common.SendMonitorMsg(common.MONITOR_LIST_SYNC_NC, "\n订单通知NC:查询未通知订单失败:---\n"+err.Error())
				time.Sleep(time.Second * 10)
				continue
			}

			if len(list) == 0 {
				time.Sleep(time.Second * 10)
				break
			}
			//并发处理
			mapreduce.MapVoid(func(source chan<- interface{}) {
				for _, item := range list {
					indexId = item.Id
					source <- item
				}
			}, func(obj interface{}) {
				item := obj.(*model.COrderSync)
				domain.OrderHandler.CSync(item)
			})
			time.Sleep(time.Second * 10)
		}
	}
}
`
	TomlTemplate = `[service]
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
    "createPreOrder",
    "getPreOrder"
]


[job]
StatisticDataByDay = "30 1 0 * * *"
LoopCSync = "30 1 0 * * *"

[cmd]
[cmd.createOrder]
flag = ["LoopOrder","LoopInvoice"]


[mysql]
[mysql.order]
tables = ["order","goods"]
[mysql.wallet]
tables = ["record"]

[redis]
[redis.order]
`
	// PkgMap imports head options. import包含选项
	PkgMap = map[string]string{
		"stirng":     `"string"`,
		"time.Time":  `"time"`,
		"gorm.Model": `"github.com/jinzhu/gorm"`,
		"fmt":        `"fmt"`,
	}
)
