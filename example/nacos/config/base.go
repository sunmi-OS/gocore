package config


var baseConfig = `
[system]
# RPC监听地址
RpcServiceHost = ""
# RPC监听端口
RpcServicePort = "8080"
# RPC-Gateway监听地址
RpcGatewayServiceHost = ""
# RPC-Gateway监听端口
RpcGatewayServicePort = "8081"

[redisDB]
remote_control = 19

[mdm_machine_online_status_info]
batchTimeout = 1000
`
