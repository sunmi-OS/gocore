目录结构
---

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