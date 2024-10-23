package middleware

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// TraceId is a middleware that injects a trace ID into the context of each request.
func TraceId(traceKey string, traceHeaders ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		md, ok := metadata.FromIncomingContext(c.Request.Context())
		if !ok {
			md = metadata.Pairs()
		}
		// traceId已存在，则复用
		if len(md.Get(traceKey)) != 0 {
			c.Next()
			return
		}
		// 取traceId，优先级按照headers顺序
		// 如：[x-b3-traceid,x-request-id]，若x-b3-traceid取不到，则取x-request-id
		for _, h := range traceHeaders {
			if v := c.GetHeader(h); len(v) != 0 {
				// 设置traceId
				md.Set(traceKey, v)
				ctx := metadata.NewIncomingContext(c.Request.Context(), md)
				c.Request = c.Request.WithContext(ctx)
				break
			}
		}
		c.Next()
	}
}
