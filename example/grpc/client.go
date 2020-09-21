package main

import (
	"github.com/sunmi-OS/gocore/rpcx"
	"github.com/sunmi-OS/gocore/xlog"
)

func main() {
	client, err := rpcx.NewGrpcClient("server_name", ":2233", nil)
	if err != nil {
		xlog.Errorf("rpcx.NewGrpcClient, err:%+v", err)
		return
	}
	grpcClient, ok := client.Next()
	if !ok {
		xlog.Error("not ok")
	}
	xlog.Debug("client ok:", grpcClient)
}
