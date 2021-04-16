package clientinterceptors

import (
	"context"
	"fmt"
	"time"

	"github.com/sunmi-OS/gocore/xlog"
	"google.golang.org/grpc"
)

var slowThreshold = (time.Millisecond * 500).Milliseconds()

// @desc 请求失败或者慢请求会打印日志
// @auth liuguoqiang 2020-04-21
// @param
// @return
func DurationInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serverName := fmt.Sprintf("host: %s, method: %s", cc.Target(), method)
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	duration := time.Since(start)
	if err != nil {
		xlog.Zap().Errorf("rpc-client-fail: server: {%s}, req: {%+v}, reply: {%+v}, err: %+v", serverName, req, reply, err)
		return err
	}
	if duration.Milliseconds() > slowThreshold {
		xlog.Zap().Warnf("rpc-client-slow: server: {%s}, duration: %dms, req: {%+v}, reply: {%+v}", serverName, duration.Milliseconds(), req, reply)
	}
	return nil
}
