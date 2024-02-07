<%: func FromREADME(buffer *bytes.Buffer) %>
<a href="https://sunmi.com"><img height="180" src="https://file.cdn.sunmi.com/gocore-logo.png"></a>

## 项目名称
> 请介绍一下你的项目吧

## 目录结构
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
> PS：使用tree工具生成目录结构 'tree -d'

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
