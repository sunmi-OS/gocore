package serverinterceptors

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StreamCrashInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})
	return handler(srv, stream)
}

func UnaryCrashInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})

		return handler(ctx, req)
	}
}

// @desc 捕获panic
// @auth liuguoqiang 2020-06-11
// @param
// @return
func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r interface{}) error {
	log.Printf("[client-fail] - %v - %s", r, debug.Stack())
	return status.Errorf(codes.Internal, "panic: %v", r)
}
