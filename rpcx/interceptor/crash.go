package interceptor

import (
	"context"
	"runtime"

	xlog2 "github.com/sunmi-OS/gocore/v2/glog/xlog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryCrash(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})
	return handler(ctx, req)
}

func StreamCrash(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})
	return handler(srv, stream)
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r interface{}) error {
	var buf [2 << 10]byte
	xlog2.Errorf("[server-panic] - %v - %s", r, string(buf[:runtime.Stack(buf[:], false)]))
	return status.Errorf(codes.Internal, "panic: %v", r)
}
