目录结构
---

通过脚手架可以快速构建项目，无需从0开始创建目录结构，无需复制历史项目统一替换项目名，通过一行命令就能快速生成项目结构和基础骨架代码，并且融入接口开发理念ADM（Api-Domain-Model）分层方式来进行开发规范的约束。

目录结构如下：

```
.
├── Dockerfile #docker镜像打包
├── README.md #项目介绍
├── app 
│   ├── api #api接口入口
│   │   └── user.go
│   ├── cronjob #定时任务,业务逻辑处理
│   │   └── sync_user.go
│   ├── def #api接口请求参数和响应参数声明
│   │   └── def.go
│   ├── domain #最终业务逻辑处理
│   │   └── user.go
│   ├── errcode #错误码声明
│   │   └── errcode.go
│   ├── job #一次性任务,常驻任务,业务逻辑
│   │   └── init_user.go
│   ├── model
│   │   └── app #model目录,app为数据库名称
│   │       ├── mysql_client.go #数据库db建表,获取db实例
│   │       └── user.go #model文件,user是表名称
│   └── routes
│       └── routers.go  #路由规则设置
├── cmd #程序启动入口
│   ├── api.go #api接口启动命令
│   ├── cron.go #定时任务启动命令
│   ├── init.go #配置,数据库,redis等初始化
│   └── job.go #一次性任务,常驻任务启动命令
├── common #存放全局变量,公共方法
├── conf #配置文件目录
│   ├── base.go #基础通用配置,最终会合并到项目配置中
│   ├── const.go #全局常量
│   └── local.go #本地配置环境变量为local时会使用
├── go.mod
├── go.sum
├── gocore.yaml #gocore项目生成配置文件
├── main.go #程序入口
└── pkg #公共工具包
```

基础结构
- main.go：作为cli主入口，无多余代码
- cmd：具体每一类型程序都需要在cmd中明确定义运行方式和cli入口
- conf：统一存放配置文件
- app：应用业务逻辑
- pkg：三方lib库存放
- common：公共变量或方法

项目目录结构遵循了一定的开发理念ADM（Api-Domain-Model），简单来说：
- Api层 称为接口服务层，负责对客户端的请求进行响应，处理接收客户端传递的参数，进行高层决策并对领域业务层进行调度，最后将处理结果返回给客户端。
- Domain层 称为领域业务层，负责对领域业务的规则处理，重点关注对数据的逻辑处理、转换和加工，封装并体现特定领域业务的规则。
- Model层 称为数据模型层，负责技术层面上对数据信息的提取、存储、更新和删除等操作，数据可来自内存，也可以来自持久化存储媒介，甚至可以是来自外部第三方系统。
