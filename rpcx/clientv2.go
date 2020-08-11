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
		Host                 string
		Name                 []string
		Timeout              string
		MaxAttempts          string
		InitialBackoff       string
		MaxBackoff           string
		BackoffMultiplier    string
		RetryableStatusCodes []string
		WaitForReady         string
		MaxTokens            string
		TokenRatio           string
		DialOptions          []grpc.DialOption
	}
)

// grpc 重试策略：重试、对冲
// 最多执行4次 RPC 请求，一个原始请求，三个重试请求，并且只有状态码为 `UNAVAILABLE` 时才重试
// MaxAttempts 最大尝试次数
// InitialBackoff 第一次重试的时间间隔
// MaxBackoff 第 n 次的重试间隔
// BackoffMultiplier 用于计算 MaxBackoff，必须大于0
// RetryableStatusCodes 匹配返回的状态码，从而进行重试
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
