package serverinterceptors

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryStatInterceptor
// 链路用时打印 and Panic拦截打印
func UnaryStatInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})
		return handler(ctx, req)
	}
}
