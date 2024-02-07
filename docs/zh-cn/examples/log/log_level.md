# glog 日志级别设置

在测试和开发环境，代码中一般会打印比较多的 INFO 及以下级别的日志来满足快速调试的需求。然而在生产环境环中，一般只允许打印 ERROR 及以上级别的日志，但又希望保留打印 INFO 级别日志的代码，以便在关键时刻可以调整日志级别到 INFO，来获取到更详细的日志数据用于定位问题。本文介绍下 gocore 项目（使用最新的 gocore 脚手架生成的项目已经自动集成）如何设置日志级别。

## 在配置文件中配置日志级别 

在 config 文件中的 base 节点下面添加日志级别配置，字段名为 logLevel，示例代码如下：

```
[base]
logLevel = "error"
```

日志级别从低到高可设置的值如下：
- debug
- info
- warn
- error
- fatal

## 添加日志初始化函数

在 /app/cmd/init.go 文件中添加如下函数：

```go
// 日志初始化函数
func initLog() {
    zap.SetLogLevel(viper.GetEnvConfig("base.logLevel").String())
}
```

代码中需要导入 github.com/sunmi-OS/gocore/v2/glog/zap 包

## 调用日志初始化函数

在 /app/cmd/api.go 文件中的 RunApi 函数中调用日志初始化函数：

```go
func RunApi(c *cli.Context) error {
    initConf()
    // 初始化日志，一定要放在初始化配置后面
    initLog()
    // 其他逻辑
    // ...
}
```

## 如何在配置中心调整日志级别后即时生效

在配置中心调整日志级别，如果要即时生效，需要添加监听配置文件变化的代码。在 /app/cmd/init.go 文件中的initConf 函数中，添加对应的监听代码，示例如下：

```
func initConf() {
	switch utils.GetRunTime() {
	case "local":
		nacos.SetLocalConfig(conf.LocalConfig)
	default:
		nacos.NewNacosEnv()
	}
	vt := nacos.GetViper()
	vt.SetBaseConfig(conf.BaseConfig)
	vt.SetDataIds(conf.ProjectName, "config", "mysql", "redis", "rocketmq")

	// 监听配置文件的变化，实时调整日志级别
	vt.SetCallBackFunc(conf.ProjectName, "config", func(namespace, group, dataId, data string) {
		initLog()
	})

    vt.NacosToViper()
}
```

## 生产环境为什么只允许打印 error 及以上级别的日志

- 在生产环境中，代码功能比较稳定，ERROR 级别的日志足以记录那些需要立即关注的问题，如异常和系统错误。
- INFO 及以下级别的日志量比较大，特别是在高流量的系统中，可能会对性能和稳定性产生负面影响，也会增加很大的存储成本。
- 生产环境中通常会有监控和警报，这些系统会触发对错误和关键事件的响应。如果日志中充斥着不必要的信息，可能会影响监控系统的准确性。

## 什么是错误日志？

错误日志是记录应用程序、系统或服务中出现的错误的日志。错误通常指的是那些导致程序无法正常运行、产生不正确结果或完全崩溃的问题。错误日志对于开发者和系统管理员来说是诊断问题的关键资源。

- 调用其他服务（API、RPC等））失败，应记录错误日志，包括错误详情和失败的上下文。
- 数据库查询、更新、写入失败。
- 文件读写、网络I/O或其他I/O操作失败时。
- 应用程序无法正确加载或解析配置文件。