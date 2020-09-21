package main

import (
	"context"

	"github.com/sunmi-OS/gocore/example/grpc/proto"
	"github.com/sunmi-OS/gocore/rpcx"
	"google.golang.org/grpc"
)

type Print struct {
}

func (p *Print) PrintOK(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	return &proto.Response{Code: 1, Data: "ok", Message: "Hello "}, nil
}

func main() {

	c := &rpcx.GrpcServerConfig{Timeout: 500}

	s := new(Print)

	rpcx.NewGrpcServer("server_name", ":2233", c).
		RegisterService(func(server *grpc.Server) {
			// register service
			proto.RegisterPrintServiceServer(server, s)
		}).Start()
}
