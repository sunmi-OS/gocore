package rpcx

import (
	"log"
	"net"
	"time"

	interceptor2 "github.com/sunmi-OS/gocore/v2/rpcx/interceptor"

	"google.golang.org/grpc"
)

type GrpcServer struct {
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

// NewGrpcServer new grpc server
func NewGrpcServer(name, addr string, cfg *GrpcServerConfig) *GrpcServer {
	server := &GrpcServer{
		Name: name,
		addr: addr,
		cfg:  cfg,
	}
	if server.cfg == nil {
		server.cfg = &GrpcServerConfig{
			Timeout: 500 * time.Millisecond,
		}
	}
	return server
}

// RegisterService .Start() 之前，必须先处理此方法
func (s *GrpcServer) RegisterService(register RegisterFn) *GrpcServer {
	s.register = register
	if s.cfg.Timeout > 0 {
		s.AddUnaryInterceptors(interceptor2.UnaryTimeout(s.cfg.Timeout))
	}
	s.AddUnaryInterceptors(interceptor2.UnaryCrash)
	s.AddStreamInterceptors(interceptor2.StreamCrash)

	options := append(s.options, interceptor2.WithUnaryServer(s.unaryInterceptors...), interceptor2.WithStreamServer(s.streamInterceptors...))
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
