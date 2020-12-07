package main

import (
	"github.com/sunmi-OS/gocore/monitor"
)

func main() {
	dingTalk := monitor.NewDingTalk("").
		WithAtMobiles([]string{}).
		WithIsAtAll(true)
	obj := monitor.NewMonitor(dingTalk)
	obj.SendTextMsg("软中台业务:代码测试请忽略")
	obj.Close(10)
}
