# 初始化&路由

初始化默认开启项：
- Gzip压缩：[gin-contrib/gzip](https://github.com/gin-contrib/gzip)
- 优雅重启&关闭：[fvbock/endless](https://github.com/fvbock/endless)
- 默认开启Recovery、Logger
- 根据utils.IsRelease设置输出
- 程序结束Close中间件连接

```go
func RunApi(c *cli.Context) error {
	defer closes.Close()

	initConf()
	initDB()
	initCache()

	if utils.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// 注册路由
	routes.Routes(r)

	err := endless.ListenAndServe(viper.C.GetString("network.ApiServiceHost")+":"+viper.C.GetString("network.ApiServicePort"), r)
	if err != nil {
		return err
	}
	return nil
}
```


默认路由
- /项目名/模块名/接口名
    - /app/user/get_user_info

考虑多入口和版本控制推荐使用方式
- /入口名称/版本号/项目名/模块名/接口名
    - /applet/v1/app/user/get_user_info

> 增加根路由监听防止被扫异常以及启动监控检查

```go

func Routes(router *gin.Engine) {

	// 根目录健康检查
	router.Any("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome GoCore Service")
	})

	user := router.Group("/app/user")
	user.Post("/get_user_info", api.GetUserInfo) //获取用户信息

}

```


