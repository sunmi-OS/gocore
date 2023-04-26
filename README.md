![logo](https://file.cdn.sunmi.com/logo.png?x-oss-process=image/resize,h_200)

介绍
---

[![Go Report Card](https://goreportcard.com/badge/github.com/sunmi-OS/gocore)](https://goreportcard.com/report/github.com/sunmi-OS/gocore)
[![GoDoc](https://godoc.org/github.com/sunmi-OS/gocore/v2?status.svg)](https://pkg.go.dev/github.com/sunmi-OS/gocore/v2)
[![Release](https://img.shields.io/github/v/release/sunmi-OS/gocore.svg?style=flat-square)](https://github.com/sunmi-OS/gocore/releases)

gocore是一款高度集成的开发框架和脚手架，支持api、rpc、job、task等开发方式，并集成各类主流开源库和中间件融入最佳实践，简化研发流程、提高效率、统一规范。

![cli](https://file.cdn.sunmi.com/gocore_cli.svg)

## 特性

- 底层基于主流框架gin、gorm、viper、zap等进行封装整合
- 提供脚手架gocore工具快速初始化项目结构、接口参数路由、数据库模型（包含逆向生成status）
- 支持多环境多套配置文件并且和nacos配置中心打通，支持热更新等特性
- 提供签名、加密、文件、邮件、随机数、链路追踪、时间、日志等基础工具
- 无侵入式理念让开发精力集中在业务层
- 通过Docker、K8S、istio等体系下建立的研发流程环境管理策略
- 封装常规阿里云中间件SLS、RocketMQ、nacos
- 开箱即用

## 安装

- 环境要求
    - Golang > 1.16
    - [Go module](https://github.com/golang/go/wiki/Modules)


### 获取项目包

```shell
> go get -u github.com/sunmi-OS/gocore/v2
```

* 脚手架安装
```shell
> go install github.com/sunmi-OS/gocore/v2/tools/gocore@latest

> gocore --version

   __ _    ___     ___    ___    _ __    ___
  / _` |  / _ \   / __|  / _ \  | '__|  / _ \
 | (_| | | (_) | | (__  | (_) | | |    |  __/
  \__, |  \___/   \___|  \___/  |_|     \___|
  |___/

gocore version v1.0.0
```


## 快速开始

创建一个示例项目
```shell
# 创建工程文件夹
> mkdir test
> cd test

# 创建yaml配置文件模板gocore.yaml
> gocore yaml create 
...
Welcome to GoCore, Configuration file has been generated.

# 修改gocore.yaml模板之后,根据yaml文件创建工程项目
> gocore service create 

   __ _    ___     ___    ___    _ __    ___
  / _` |  / _ \   / __|  / _ \  | '__|  / _ \
 | (_| | | (_) | | (__  | (_) | | |    |  __/
  \__, |  \___/   \___|  \___/  |_|     \___|
  |___/

Run go mod init.
[11/11] Initialize the Request return parameters... 100% [========================================]   
Run go mod tidy .
Run go fmt .
goimports -l -w .
Welcome to GoCore, the project has been initialized.

# 下次迭代增加新的接口或数据表更新代码
> gocore service create 

```

工程创建时导入已有数据库
```shell
# 创建工程文件夹
> mkdir test 
> cd test

# 创建yaml配置文件模板gocore.yaml
> gocore yaml create 

# 创建连接数据库的配置文件模板mysql.yaml
> gocore mysql create_yaml 

# 修改mysql.yaml之后,连接数据库将字段合并到gocore.yaml
> gocore mysql add 

# 修改gocore.yaml模板之后,根据yaml文件创建工程项目
> gocore service create 
```


## 配置文件

```yaml
service:
  projectName: demo #项目名称
  version: v1.0.0 #项目版本号
config:
  cNacos: true #是否使用nacos
  cRocketMQConfig: true #是否使用rocketMQ
  cMysql: #mysql配置
    - name: app #数据库名称
      hotUpdate: false #是否热更新
      models: #model文件
        - name: user #表名称
          auto: false #是否自动建表
          fields: #表字段,gorm规则,一行一个自动
            - column:id;primary_key;type:int AUTO_INCREMENT
            - column:name;type:varchar(100) NOT NULL;default:'';comment:'用户名';unique_index
          comment: 用户表 #表备注
  cRedis: #redis配置
    - name: default #redis名称
      hotUpdate: false #是否热更新
      index:
        db0: 0 #选择第几个db
rpcEnable: false #是否生成rpc服务层
httpApiEnable: true #是否生成接口程序
jobEnable: true #是否生成常驻任务
httpApis:
  host: 0.0.0.0 #api接口监听ip地址
  port: "80" #api接口监听ip端口
  apis:
    - prefix: /app/user #api接口前缀
      moduleName: user #模块名称
      handle: #api接口
        - name: GetUserInfo #api接口方法名称,完整路由是/app/user/GetUserInfo
          method: Any
          requestParams: #api接口请求参数
            - name: uid #字段名称
              type: int #字段类型
              comment: 用户ID #字段备注
              validate: required,min=1,max=100000 #validate校验规则
          responseParams: #api响应参数
            - name: detail  #字段名称
              type: '*User'  #字段类型,非基础字段类型,表示嵌套结构体,引用params中的结构体
              comment: 用户详情 #字段备注
              validate: ""
            - name: list
              type: '[]*User'
              comment: 用户列表
              validate: ""
          comment: 获取用户信息
  params:
    User:
      - name: uid
        type: int
        comment: 用户ID
        validate: ""
      - name: name
        type: string
        comment: 用户名
        validate: ""
jobs:
  - name: InitUser #一次性任务,常驻任务方法名称
    comment: 初始化默认用户 #一次性任务,常驻任务备注
```

## 生成的工程目录结构
使用三层架构(http服务: api->biz->dal,rpc服务: rpc->biz->dal)：

- api(rpc) 接口表示层 Application Programming Interface
  - 定义接口名称，校验入参，调用biz层方法处理业务逻辑并返回响应数据
  - 只能调用biz层方法，禁止调用dal层方法
- biz 业务逻辑层 Business Logic Layer
  - 业务逻辑处理层，接收api层传入的参数结合调用dal层方法完成业务逻辑处理并返回必要数据
  - 禁止调用api层方法
- dal 数据访问层 Data Access Layer
  - 负责对DB的访问，本层禁止相互的方法调用
  - 禁止调用api和biz层方法

目录结构说明：
```
├── app                  // 源代码
│  ├── api               // 接口表示层，无http服务的话删除此文件夹
│  ├── rpc               // rpc服务表示层，无rpc服务删除此文件夹
│  ├── biz               // 业务逻辑层
│  ├── dal               // 数据访问层
│  ├── middleware        // 中间件
│  ├── cmd               // 任务启动入口和定义各组件初始化方法
│  │  ├── api.go
│  │  ├── init.go
│  │  └── job.go
│  ├── conf              // 配置文件
│  │  ├── base.go        // 基本配置
│  │  └── local.go       // 用于本地调试配置文件，本地环境变量需要设置RUN_TIME=local
│  ├── errcode           // 错误和错误码定义
│  │  └── errcode.go
│  ├── job               // 任务定义入口，定时任务、一次性任务、消费队列任务
│  ├── param             // 入参和出参结构体定义
│  │  └── user.go
│  ├── pkg               // 依赖包和三方包
│  │  ├── locationtools  // 三方包封装示例
│  │  │  └── country.go
│  │  └── util           // 实现的常用方法
│  │      └── util.go
│  ├── route             // 路由定义
│  │   └── routers
│  ├── go.mod
│  ├── go.mod
│  └── main.go           // 入口文件
├── .gitignore
├── CODEOWNERS           // 用来定义谁负责仓库中的特定文件或目录
├── Dockerfile
└── README.md
```