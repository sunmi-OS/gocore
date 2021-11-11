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
```bash
# 创建工程文件夹
> mkdir test
> cd test

# 创建yaml配置文件模板gocore.yaml
> gocore conf create 
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
```bash
# 创建工程文件夹
> mkdir test 
> cd test

# 创建yaml配置文件模板gocore.yaml
> gocore conf create 

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
  cNacos:
    env: false #是否使用环境变量
    rocketMQConfig: true #是否使用rocketMQ
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
nacosEnable: true #是否使用nacos
httpApiEnable: true #是否生成接口程序
cronJobEnable: true #是否生成定时任务
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
              required: true #是否必填
              type: int #字段类型
              comment: 用户ID #字段备注
              validate: required,min=1,max=100000 #validate校验规则
          responseParams: #api响应参数
            - name: detail  #字段名称
              required: true #是否必填
              type: '*User'  #字段类型,非基础字段类型,表示嵌套结构体,引用params中的结构体
              comment: 用户详情 #字段备注
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
  - spec: '@every 30m' #定时任务规则,参考:github.com/robfig/cron
    job:
      name: SyncUser #定时任务方法名称
      comment: 同步用户 #定时任务备注
jobs:
  - name: InitUser #一次性任务,常驻任务方法名称
    comment: 初始化默认用户 #一次性任务,常驻任务备注
```


## 联系我们

欢迎加入`gocore`QQ群：1004023331 一起沟通讨论

![qq](https://file.cdn.sunmi.com/qq.png?x-oss-process=image/resize,h_200)
