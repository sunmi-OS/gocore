package main

import (
	"context"
	"time"

	"github.com/sunmi-OS/gocore/example/grpc/proto"
	"github.com/sunmi-OS/gocore/rpcx"
	"google.golang.org/grpc"
)

type Print struct {
}

func (p *Print) PrintOK(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	time.Sleep(1000 * time.Millisecond)
	return &proto.Response{Code: 1, Data: "ok", Message: "Hello client"}, nil
}

func main() {

	c := &rpcx.GrpcServerConfig{Timeout: 500 * time.Millisecond}

	s := new(Print)

	rpcx.NewGrpcServer("server_name", ":2233", c).
		RegisterService(func(server *grpc.Server) {
			// register service
			proto.RegisterPrintServiceServer(server, s)
		}).Start()
}
