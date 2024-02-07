# sls阿里云日志库

- 阿里云日志输出需要配置阿里云连接信息

```
[alilog]
Project = "" #项目名称
Endpoint = ""  #连接地址
AccessKey = ""
SecretKey = ""
```

- 代码示例

```go
package main

import (
	"errors"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
    "github.com/sunmi-OS/gocore/v2/glog/sls"
)

func main() {
	s := struct {
		Name string
		Age  int
	}{
		Name: "Jerry",
		Age:  18,
	}

	conf := `
[alilog]
Project = ""
Endpoint = ""
AccessKey = ""
SecretKey = ""
	`
	viper.MergeConfigToToml(conf)
	sls.InitLog("alilog", "message-core")
	sls.SetGLog()

	glog.InfoF("%+v", s)
	glog.Debug("zap debug")
	glog.Warn("zap warn")
	glog.Error("zap error", "呵呵")
	glog.ErrorF("s.dao.PartnerById(%d),err:%+v", 10086, errors.New("不存在此id"))
	glog.ErrorF("s.dao.CreateOrder(%+v),err:%+v", s, errors.New("创建订单失败"))
}
```

```go
// 将glog新增输出到阿里云
func SetGLog() {
	glog.SetLogger("alilog", &sls.LogClient)
}
```