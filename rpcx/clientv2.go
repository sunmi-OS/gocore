package rpcx

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sunmi-OS/gocore/rpcx/clientinterceptors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

type (
	ClientConfig struct {
		Host                 string            //RPC链接
		Name                 []string          //RPC包名.服务名 如shop_base.ShopBase
		Timeout              string            //超时时间,默认:5s
		MaxAttempts          string            //最大重试次数,必须是大于 1 的整数，对于大于5的值会被视为5,默认:4
		InitialBackoff       string            //第一次重试默认时间间隔,必须具有大于0,默认:2s
		MaxBackoff           string            //最大重试时间间隔,必须具有大于0,默认:3s
		BackoffMultiplier    string            //间隔增量乘数因子,大于零,默认:1
		RetryableStatusCodes []string          //重试会根据请求返回的状态码是否符合 retryableStatusCodes来进行重试请求,默认:UNAVAILABLE
		WaitForReady         string            //如果为false，则RPC将在连接到服务器的瞬间失败时立即中止。否则，gRPC会尝试连接，直到超过截止日期。默认:true
		MaxTokens            string            //如果 token_count <= ( maxTokens / 2), 则关闭重试策略，直到 token_count > (maxTokens/2)，恢复重试,默认10
		TokenRatio           string            //成功 RPC 将会递增 token_count * tokenRatio
		DialOptions          []grpc.DialOption //RPC链接可选参数
	}
)

// @desc RPC客户端初始化链接配置
// @auth liuguoqiang 2020-08-12
// @param
// @return
func NewClientConfig(host string, name []string) (*ClientConfig, error) {
	if host == "" {
		return nil, fmt.Errorf("rpc链接地址不能为空")
	}
	if len(name) == 0 {
		return nil, fmt.Errorf("服务名不能为空")
	}
	return &ClientConfig{
		Host:                 host,
		Name:                 name,
		Timeout:              "5s",
		MaxAttempts:          "2",
		InitialBackoff:       "2s",
		MaxBackoff:           "3s",
		BackoffMultiplier:    "1",
		RetryableStatusCodes: []string{"UNAVAILABLE"},
		WaitForReady:         "true",
		MaxTokens:            "10",
		TokenRatio:           "0.1",
	}, nil
}

// @desc 设置默认配置
// @auth liuguoqiang 2020-08-12
// @param
// @return
func (clientConfig *ClientConfig) WithDefaultServiceConfig() grpc.DialOption {
	nameArr := make([]map[string]string, 0)
	for k1 := range clientConfig.Name {
		name := map[string]string{
			"service": clientConfig.Name[k1],
		}
		nameArr = append(nameArr, name)
	}
	nameByte, err := json.Marshal(nameArr)
	if err != nil {
		log.Print(string(nameByte))
	}

	retryableStatusCodesArr := make([]string, 0)
	for k1 := range clientConfig.RetryableStatusCodes {
		retryableStatusCodesArr = append(retryableStatusCodesArr, clientConfig.RetryableStatusCodes[k1])
	}
	retryableStatusCodesByte, err := json.Marshal(retryableStatusCodesArr)
	if err != nil {
		log.Print(string(retryableStatusCodesByte))
	}
	retryPolicy := fmt.Sprintf(`{
    "LoadBalancingPolicy":"%s",
    "loadBalancingConfig":[
        {
            "round_robin":{
            }
        }
    ],
    "methodConfig":[
        {
            "name":%s,
            "waitForReady":%s,
            "timeout":"%s",
            "retryPolicy":{
                "MaxAttempts":%s,
                "InitialBackoff":"%s",
                "MaxBackoff":"%s",
                "BackoffMultiplier":%s,
                "RetryableStatusCodes":%s
            }
        }
    ],
    "retryThrottling":{
        "maxTokens":%s,
        "tokenRatio":%s
    }
}`, roundrobin.Name, string(nameByte), clientConfig.WaitForReady, clientConfig.Timeout, clientConfig.MaxAttempts, clientConfig.InitialBackoff, clientConfig.MaxBackoff, clientConfig.BackoffMultiplier, string(retryableStatusCodesByte), clientConfig.MaxTokens, clientConfig.TokenRatio)
	return grpc.WithDefaultServiceConfig(retryPolicy)
}

func (clientConfig *ClientConfig) buildDialOptions() []grpc.DialOption {
	interceptor := make([]grpc.UnaryClientInterceptor, 0)
	interceptor = append(interceptor, clientinterceptors.DurationInterceptor)
	options := []grpc.DialOption{
		grpc.WithInsecure(),
		WithUnaryClientInterceptors(interceptor...),
		clientConfig.WithDefaultServiceConfig(),
	}

	return append(options, clientConfig.DialOptions...)
}
