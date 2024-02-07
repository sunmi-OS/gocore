# mongodb

使用 mongodb 官方driver  go.mongodb.org/mongo-driver/mongo

## 使用方法

### 修改配置信息

在配置文件配置mongodb信息，假如库名是test

```
[mongoTest]
Endpoint = ["xxx.xxx.com:3717","xxx2.xxx.com:3717"]    #数据库连接地址
Name = "test"           #数据库名称
User = "xxx"            #数据库用户名
Passwd = "xxxxxxx"      #数据库密码
ReplicaSet = "xxxx" # 副本集，副本集模式集群必须填写，分片集群和单节点无需此字段
MaxPoolSize = 20        #最大连接数
MinPoolSize = 10        #最小连接数
```

### 初始化

在 app/cmd/init.go 添加初始化方法

```go
func initMogoDB() {
	mongodb.NewDB("mongoTest")
}
```

在服务启动方法里面调用初始化方法（例如在/app/cmd/api.go 文件的 RunApi方法）

```go
func RunApi(c *cli.Context) error {
	defer closes.Close()

	initConf()
	...
	initMogoDB()
  
  // ...
  
}
```

### 操作 mongodb

在 app/dal 新建对应 mongodb 数据库文件，例如 test，在 test 文件下创建 client.go文件，在次文件创建获取 db 实例方法 Orm()

```go
func Orm() *mongo.Database {
    return mongodb.GetDB("mongoTest")
}
```