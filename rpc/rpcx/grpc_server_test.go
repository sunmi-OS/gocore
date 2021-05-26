package rpcx

import (
	"testing"
)

const (
	MaxMessageSize = 1024 * 1024 * 64 // 64MB
)

func TestNewGrpcServer(t *testing.T) {
	//addr := ""
	//xlog.Info("Grpc Server:", addr)
	//
	//// start grpc server
	//server := NewGrpcServer("server_name", addr, &GrpcServerConfig{Timeout: 0})
	//ops := []grpc.ServerOption{
	//	grpc.MaxRecvMsgSize(MaxMessageSize),
	//	grpc.MaxSendMsgSize(MaxMessageSize),
	//}
	//server.AddOptions(ops...)
	//
	//server.RegisterService(func(server *grpc.Server) {
	//	// todo: register your grpc server
	//
	//})
	//go func() {
	//	server.Start()
	//}()
	//time.Sleep(time.Hour)
}
