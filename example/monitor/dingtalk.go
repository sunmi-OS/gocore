package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/monitor"
)

func main() {
	resp, err := monitor.NewDingTalk("url").
		WithAtMobiles([]string{}).
		WithIsAtAll(true).
		SendTextMsg("软中台业务:代码测试请忽略")
	fmt.Printf("%#v\n", string(resp))
	fmt.Printf("%#v\n", err)
}
