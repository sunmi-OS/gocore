![logo](https://file.cdn.sunmi.com/logo.png?x-oss-process=image/resize,h_200)

介绍
---

gocore是一款高度集成的开发框架和脚手架，支持api、rpc、job、task等开发方式，并且集成各类主流开源库和中间件融入最佳实践，最终实现简化流程、提高效率、统一规范。


## 安装

- 环境要求
  - Golang > 1.16
  - [Go module](https://github.com/golang/go/wiki/Modules)
  

### 获取项目包

```

> go get -u github.com/sunmi-OS/gocore/v2

```

* 脚手架安装
```

> go get -u github.com/sunmi-OS/gocore/v2/tool/gocore

> gocore --version

   __ _    ___     ___    ___    _ __    ___
  / _` |  / _ \   / __|  / _ \  | '__|  / _ \
 | (_| | | (_) | | (__  | (_) | | |    |  __/
  \__, |  \___/   \___|  \___/  |_|     \___|
  |___/

gocore version v1.0.0

```

## 配置文件

```yaml
service:
  projectName: demo
  version: v1.0.0
config:
  cNacos:
    env: false
    rocketMQConfig: true
  cMysql:
  - name: app
    hotUpdate: false
    models:
    - name: user
      auto: false
      fields:
      - name: ""
        gormRule: column:id;primary_key;type:int AUTO_INCREMENT
      - name: ""
        gormRule: column:name;type:varchar(100) NOT NULL;default:'';comment:'用户名';unique_index
      comment: 用户表
  cRedis:
  - name: default
    hotUpdate: false
    index:
      db0: 0
nacosEnable: true
httpApiEnable: true
cronJobEnable: true
jobEnable: true
httpApis:
  host: 0.0.0.0
  port: "80"
  apis:
  - prefix: /app/user
    moduleName: user
    handle:
    - name: GetUserInfo
      method: Any
      requestParams:
      - name: uid
        required: true
        type: int
        comment: 用户ID
        validate: required,min=1,max=100000
      responseParams:
      - name: detail
        required: true
        type: '*User'
        comment: 用户详情
        validate: ""
      - name: list
        required: true
        type: '[]*User'
        comment: 用户列表
        validate: ""
      comment: 获取用户信息
  params:
    User:
    - name: uid
      required: true
      type: int
      comment: 用户ID
      validate: ""
    - name: name
      required: true
      type: string
      comment: 用户名
      validate: ""
cronJobs:
- spec: '@every 30m'
  job:
    name: SyncUser
    comment: 同步用户
jobs:
- name: InitUser
  comment: 初始化默认用户
```


## 特性

* 脚手架

通过gocore工具快速初始化项目结构、接口参数路由、数据库模型

* 配置文件

支持多环境多套配置文件并且和nacos配置中心打通，支持热更新等特性

* Service Mesh支持

对于使用Istio的开发者，KT提供指向本地服务Version标签和自定义标签来精细控制流量

* 丰富的基础工具

提供签名、加密、文件、邮件、随机数、链路追踪、时间、日志等基础工具

* 最佳实践

基于Docker、K8S、istio等体系下建立的研发流程环境管理策略

## 联系我们

请加入`gocore`QQ群：1004023331

![qq](https://file.cdn.sunmi.com/qq.png?x-oss-process=image/resize,h_200)
