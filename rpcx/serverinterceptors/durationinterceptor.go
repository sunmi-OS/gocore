package serverinterceptors

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
		log.Printf("[sever-fail] %s - %s", addr, err.Error())
	} else if duration > serverSlowThreshold {
		log.Printf("[sever-slow] - %s - %s - %s", addr, method, string(content))
	} else {
		log.Printf("[sever-call] %s - %s - %s", addr, method, string(content))
	}
}
