# HTTP 初始化&路由注册

在 HTTP 接口开发过程中，考虑实际场景初始化默认开启：
- 优雅重启&关闭
- 默认开启 Recovery、Logger、Tracing
- 默认集成 Prometheus、健康检查接口、Pprof
- 根据 utils.IsRelease 设置终端输出
- 程序结束 Close 中间件连接

```go
func RunApi(c *cli.Context) error {
    initConf()
    initDB()
    
    isDebugMode := true
    if utils.IsRelease() {
    isDebugMode = false
    }
    
    gs := api.NewGinServer(
    api.WithServerDebug(isDebugMode),
    api.WithServerHost(viper.C.GetString("network.ApiServiceHost")),
    api.WithServerPort(viper.C.GetInt("network.ApiServicePort")),
    api.WithOpenTrace(true),
    )
    // init route
    route.Routes(gs.Gin)
    gs.Start()
    
    return nil
}
```