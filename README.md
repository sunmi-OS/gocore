![logo](https://file.cdn.sunmi.com/logo.png?x-oss-process=image/resize,h_200)

介绍
---

gocore是一款高度集成的开发框架和脚手架，支持api、rpc、job、task等开发方式，并且集成各类主流开源库和中间件融入最佳实践，最终实现简化流程、提高效率、统一规范。


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


## 快速开始

创建一个示例项目
```bash
# 创建工程文件夹
> mkdir test
> cd test

# 创建yaml配置文件模板gocore.yaml
> gocore conf create 
...
Welcome to GoCore, Configuration file has been generated.

# 修改gocore.yaml模板之后,根据yaml文件创建工程项目
> gocroe service create 

# 下次迭代增加新的接口或数据表更新代码
> gocroe service create 

```

工程创建时导入已有数据库
```bash
# 创建工程文件夹
> mkdir test 
> cd test

# 创建yaml配置文件模板gocore.yaml
> gocore conf create 

# 创建连接数据库的配置文件模板mysql.yaml
> gocore mysql create_yaml 

# 修改mysql.yaml之后,连接数据库将字段合并到gocore.yaml
> gocroe mysql add 

# 修改gocore.yaml模板之后,根据yaml文件创建工程项目
> gocroe service create 
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

## 联系我们

欢迎加入`gocore`QQ群：1004023331 一起沟通讨论

![qq](https://file.cdn.sunmi.com/qq.png?x-oss-process=image/resize,h_200)
