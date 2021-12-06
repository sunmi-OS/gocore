gocore.yaml配置文件介绍
---

脚手架工具基于yaml配置文件生产代码，主要特性：
- 支持Api、cronjob、job类型
  - 自动生成Api接口
    - 路由
    - 入口方法
    - 参数结构
    - bind参数验证
  - 自动生成cronjob和job入口cmd
- 配置文件
  - 集成配置中间nacos
    - 支持AK&SK秘钥鉴权
    - 支持从nacos或本地读取配置
    - 使用nacos支持热更新
  - 配置文件多环境切换
  - 配置文件融入bin包
- 中间件
  - mysql
  - redis
  - rocketMQ


```yaml
service:
  projectName: demo #项目名称
  version: v1.0.0 #项目版本号
config: 
  # @TODO 需要修改
  cNacos: false #是否使用nacos
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