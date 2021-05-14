package rpcx

import (
	"log"
	"net"
	"time"

	"github.com/sunmi-OS/gocore/rpcx/serverinterceptors"
	"google.golang.org/grpc"
)

type (
	// Deprecated
	RpcServer struct {
		*baseRpcServer

		register RegisterFn
	}

	GrpcServer struct {
		Name               string
		addr               string
		isPre              bool
		cfg                *GrpcServerConfig
		register           RegisterFn
		server             *grpc.Server
		options            []grpc.ServerOption
		streamInterceptors []grpc.StreamServerInterceptor
		unaryInterceptors  []grpc.UnaryServerInterceptor
	}
)

// Deprecated
// 推荐使用 NewGrpcServer
// @desc rpc服务端初始化入口函数
// @auth liuguoqiang 2020-06-11
// @param timeout 为0时，不做超时处理
// @return
func NewRpcServer(name, addr string, timeout int64, register RegisterFn) *RpcServer {
	server := &RpcServer{
		baseRpcServer: newBaseRpcServer(addr),
		register:      register,
	}
	setupInterceptors(server, int(timeout))
	return server
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
		//serverinterceptors.UnaryCrashInterceptor(),
		serverinterceptors.UnaryStatInterceptor(),
	}
	unaryInterceptors = append(unaryInterceptors, s.unaryInterceptors...)

	streamInterceptors := []grpc.StreamServerInterceptor{
		serverinterceptors.StreamCrashInterceptor,
	}
	streamInterceptors = append(streamInterceptors, s.streamInterceptors...)

	options := append(s.options, WithUnaryServerInterceptors(unaryInterceptors...), WithStreamServerInterceptors(streamInterceptors...))
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
func setupInterceptors(server *RpcServer, timeout int) {
	if timeout > 0 {
		server.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(time.Duration(timeout) * time.Millisecond))
	}
}

// NewGrpcServer new grpc server
func NewGrpcServer(name, addr string, cfg *GrpcServerConfig) *GrpcServer {
	server := &GrpcServer{
		Name: name,
		addr: addr,
		cfg:  cfg,
	}
	if server.cfg == nil {
		server.cfg = defaultServerConfig()
	}
	return server
}

// RegisterService .Start() 之前，必须先处理此方法
func (s *GrpcServer) RegisterService(register RegisterFn) *GrpcServer {
	s.register = register
	if s.cfg.Timeout > 0 {
		s.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(s.cfg.Timeout))
	}
	s.AddUnaryInterceptors(serverinterceptors.UnaryCrashInterceptor)
	s.AddStreamInterceptors(serverinterceptors.StreamCrashInterceptor)
	options := append(s.options, WithUnaryServerInterceptors(s.unaryInterceptors...), WithStreamServerInterceptors(s.streamInterceptors...))

	s.server = grpc.NewServer(options...)
	s.isPre = true
	return s
}

// Start start grpc server.
func (s *GrpcServer) Start() {
	if !s.isPre {
		log.Fatal("before start, you must call server.RegisterService().")
	}
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
	}
	if s.server != nil {
		s.register(s.server)
	}

	if err = s.server.Serve(lis); err != nil {
		s.server.GracefulStop()
		log.Fatal(err)
	}
}

func (s *GrpcServer) Close() {
	if s.server != nil {
		s.server.Stop()
	}
}
