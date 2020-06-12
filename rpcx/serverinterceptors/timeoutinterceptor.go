package serverinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// @desc 超时插件
// @auth liuguoqiang 2020-06-11
// @param
// @return
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
		defer cancel()
		return handler(ctx, req)
	}
}
