package main

import (
	"context"

	"github.com/sunmi-OS/gocore/example/grpc/proto"
	"github.com/sunmi-OS/gocore/rpcx"
	"github.com/sunmi-OS/gocore/xlog"
)

func main() {
	client, err := rpcx.NewGrpcClient("server_name", ":2233", nil)
	if err != nil {
		xlog.Errorf("rpcx.NewGrpcClient, err:%+v", err)
		return
	}
	conn, ok := client.Conn()
	if !ok {
		xlog.Error("not ok")
		return
	}
	printGRPC := proto.NewPrintServiceClient(conn)

	req := &proto.Request{}
	rsp, err := printGRPC.PrintOK(context.Background(), req)
	if err != nil {
		xlog.Errorf("printGRPC.PrintOK(%+v), err:%+v", req, err)
		return
	}
	xlog.Info(rsp)
}
