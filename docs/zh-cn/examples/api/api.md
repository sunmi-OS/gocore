# API接口开发

Gocore基于主流网络框架Gin来实现请求路由，并且基于官方最佳实践进行集成。


并且推荐统一使用POST方式请求和JSON方式传参，便于在签名鉴权加密统一方式。








### context

```go
ctx := api.NewContext(g)

...

type Context struct {
    *gin.Context
    C context.Context
    R Response
    T *utils.TraceHeader
}

```


