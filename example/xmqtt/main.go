package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sunmi-OS/gocore/xlog"
	"github.com/sunmi-OS/gocore/xmqtt"
)

func main() {
	// 初始化参数和配置
	mq := xmqtt.New(&xmqtt.Config{
		Host:           "host",
		Port:           2233,
		ClientId:       "ClientId",
		Uname:          "uname",
		Password:       "passwd",
		IsCleanSession: true,
	})
	// 设置Mqtt连接监听
	mq.OnConnectListener(mq.DefaultOnConnectFunc)
	// 设置Mqtt断开连接监听
	mq.OnConnectLostListener(func(client mqtt.Client, err error) {
		xlog.Warnf("mq[%+t] lost connection,err:%+v", client.IsConnected(), err)
	})
	// error
	// 启动连接
	if err := mq.StartAndConnect(); err != nil {
		xlog.Error(err)
		return
	}
}
