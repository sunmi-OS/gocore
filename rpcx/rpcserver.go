package rpcx

import (
	"log"
	"net"
	"time"

	"github.com/sunmi-OS/gocore/rpcx/serverinterceptors"
	"google.golang.org/grpc"
)

type (
	RpcServer struct {
		*baseRpcServer

		register RegisterFn
	}
)

// @desc rpc服务端初始化入口函数
// @auth liuguoqiang 2020-06-11
// @param
// @return
func NewRpcServer(listenOn string, timeout int64, register RegisterFn) (*RpcServer, error) {
	var err error
	server := &RpcServer{
		baseRpcServer: newBaseRpcServer(listenOn),
		register:      register,
	}
	if err = setupInterceptors(server, timeout); err != nil {
		return nil, err
	}
	return server, nil
}

// @desc 启动rpc
// @auth liuguoqiang 2020-06-11
// @param
// @return
func (s *RpcServer) Start() {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal(err)
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryCrashInterceptor(),
		serverinterceptors.UnaryStatInterceptor(),
	}
	unaryInterceptors = append(unaryInterceptors, s.unaryInterceptors...)
	streamInterceptors := []grpc.StreamServerInterceptor{
		serverinterceptors.StreamCrashInterceptor,
	}
	streamInterceptors = append(streamInterceptors, s.streamInterceptors...)
	options := append(s.options, WithUnaryServerInterceptors(unaryInterceptors...),
		WithStreamServerInterceptors(streamInterceptors...))
	server := grpc.NewServer(options...)

	s.register(server)
	err = server.Serve(lis)
	server.GracefulStop()
	log.Fatal(err)
}

// @desc 设置超时
// @auth liuguoqiang 2020-06-11
// @param
// @return
func setupInterceptors(server *RpcServer, timeout int64) error {
	if timeout > 0 {
		server.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
			time.Duration(timeout) * time.Millisecond))
	}
	return nil
}
