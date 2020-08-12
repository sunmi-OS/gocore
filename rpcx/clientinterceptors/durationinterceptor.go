package clientinterceptors

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/2276282419/gocore/rpcx/logx"

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
	duration := time.Since(start)
	if err != nil {
		errMsg := err.Error()
		reqBtye, err := json.Marshal(req)
		if err != nil {
			fmt.Printf("%#v\n", err)
		}
		logx.LoggerObj.Error("rpc-client-fail", map[string]string{"duration": fmt.Sprintf("%d", duration), "serverName": serverName, "req": string(reqBtye), "err": errMsg})
	} else {
		elapsed := time.Since(start)
		if elapsed > slowThreshold {
			reqBtye, err := json.Marshal(req)
			if err != nil {
				fmt.Printf("%#v\n", err)
			}
			replyBtye, err := json.Marshal(reply)
			if err != nil {
				fmt.Printf("%#v\n", err)
			}
			logx.LoggerObj.Error("rpc-client-slow", map[string]string{"elapsed": fmt.Sprintf("%d", duration), "serverName": serverName, "req": string(reqBtye), "reply": string(replyBtye)})
		}
	}

	return err
}
