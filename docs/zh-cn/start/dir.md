目录结构
---

通过脚手架可以快速构建项目，无需从0开始创建目录结构，无需复制历史项目统一替换项目名，通过一行命令就能快速生成项目结构和基础骨架代码，并且使用分层方式来进行开发规范的约束。

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

目录结构如下：

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
