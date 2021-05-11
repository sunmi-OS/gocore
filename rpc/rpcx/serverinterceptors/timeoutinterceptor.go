package serverinterceptors

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// @desc 超时插件
// @auth liuguoqiang 2020-06-11
// @param
// @return
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		done := make(chan error, 1)
		h := func() {
			resp, err = handler(ctx, req)
			if err != nil {
				done <- err
				return
			}
			done <- nil
		}
		go h()

		select {
		case err := <-done:
			if err != nil {
				return nil, err
			}
			return resp, nil
		case <-ctx.Done():
			return nil, fmt.Errorf("%s request timeout", info.FullMethod)
		}
	}
}
