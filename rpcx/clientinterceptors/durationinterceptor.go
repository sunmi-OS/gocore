package clientinterceptors

import (
	"context"
	"log"
	"path"
	"time"

	"google.golang.org/grpc"
)

const slowThreshold = time.Millisecond * 500

// @desc 请求失败或者慢请求会打印日志
// @auth liuguoqiang 2020-04-21
// @param
// @return
func DurationInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serverName := path.Join(cc.Target(), method)
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		log.Printf("[client-fail] - %v - %s - %v - %s", time.Since(start), serverName, req, err.Error())
	} else {
		elapsed := time.Since(start)
		if elapsed > slowThreshold {
			log.Printf("[client-slow] ok - slowcall-%v - %s - %v - %v", elapsed, serverName, req, reply)
		}
	}

	return err
}
