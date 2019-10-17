package config

var localConfig = `
[remotemanageDB]
dbHost = "dev.db.sunmi.com"                #数据库连接地址
dbName = "remote_manage"                 #数据库名称
dbUser = "kingshardadmin"           #数据库用户名
dbPasswd = "Kwbd7246005c039789d9"   #数据库密码
dbPort = "3306"                     #数据库端口号
dbOpenconns_max = 20                #最大连接数
dbIdleconns_max = 20                #最大空闲连接
dbType = "mysql"

[redisServer]
host = "redis.database"
port = ":6379"
auth = "devsunmiredis666"
prefix = "tob_"
encryption = 1

[kafkaClient]
topic = "mdm_machine_online_status_info"
brokers = ["192.168.3.217:9092","192.168.3.218:9092","192.168.3.219:9092"]
`
