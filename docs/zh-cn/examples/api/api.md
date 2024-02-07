# API接口开发

gocore 框架最核心的能力是围绕着日常业务开发设计的，其中以 Api 为最高频的场景，根据在项目研发过程中的积累的经验进行了封装和集成。

特性：
- 基于主流框架 [Gin](https://github.com/gin-gonic/gin)
- 统一 POST 请求和 JSON 传参（统一签名鉴权加密）
- 统一返回数据格式（Code、Data、Msg）
- 集成 context.Context 包，自动捕获链路信息

## 最佳实践

核心理念：让研发更多关注业务，通过非侵入式来提供各项扩展能力，所以不使用各类侵入式的微服务框架:
- 基于 K8S+Istio 体系来实现微服务中的服务发现、链路追踪、流量治理、熔断降级等机制，并且能够实现多环境隔离等功能
- 基于 Kong+plugin 来代理入口流量，完成统一的日志记录、签名、鉴权和加密等统一能力

![gocore-01.jpg](https://file.cdn.sunmi.com/gocore/gocore-01.jpg)

### context

通过 **api.NewContext(g)** 获取一个新的 context，标准 gin 提供的 context 正常使用，并且集成了：
- 官方 context.Context 包（Deadline、Value、Done）
- 标准返回结构体
- 链路追踪信息

```go
// GetUserInfo 获取用户信息
func GetUserInfo(g *gin.Context) {
    ctx := api.NewContext(g)

...

type Context struct {
    *gin.Context
    C context.Context
    R Response
    T *utils.TraceHeader
}
```


