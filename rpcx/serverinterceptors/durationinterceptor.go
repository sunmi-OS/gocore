package serverinterceptors

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sunmi-OS/gocore/rpcx/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const serverSlowThreshold = time.Millisecond * 500

func UnaryStatInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			logDuration(ctx, info.FullMethod, req, duration)
		}()

		return handler(ctx, req)
	}
}

// @desc 打印请求时间
// @auth liuguoqiang 2020-06-11
// @param
// @return
func logDuration(ctx context.Context, method string, req interface{}, duration time.Duration) {
	var addr string
	client, ok := peer.FromContext(ctx)
	if ok {
		addr = client.Addr.String()
	}
	content, err := json.Marshal(req)
	if err != nil {
		logx.LoggerObj.Error("rpc-sever-fail", map[string]string{"addr": addr, "method": method, "content": err.Error(), "duration": fmt.Sprintf("%d", duration/time.Millisecond)})
	} else if duration > serverSlowThreshold {
		logx.LoggerObj.Info("rpc-sever-slow", map[string]string{"addr": addr, "method": method, "content": string(content), "duration": fmt.Sprintf("%d", duration/time.Millisecond)})
	} else {
		logx.LoggerObj.Info("rpc-sever-call", map[string]string{"addr": addr, "method": method, "content": string(content), "duration": fmt.Sprintf("%d", duration/time.Millisecond)})
	}
}
