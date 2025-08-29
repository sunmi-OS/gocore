package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sunmi-OS/gocore/v2/utils"
	"google.golang.org/grpc/metadata"
)

// TraceId is a middleware that injects a trace ID into the context of each request.
func TraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		md, ok := metadata.FromIncomingContext(c.Request.Context())
		if !ok {
			md = metadata.Pairs()
		}
		// traceId已存在，则复用
		if len(md.Get(utils.XB3TraceId)) > 0 {
			c.Next()
			return
		}
		// 去header取traceId
		traceId := c.GetHeader(utils.XB3TraceId)
		// 找不到x-b3-traceid，用x-request-id
		if traceId == "" {
			traceId = c.GetHeader(utils.XRequestId)
		}
		// if the traceId is absent, generate it
		if traceId == "" {
			traceId = newTraceId()
		}
		// 设置traceId
		md.Set(utils.XB3TraceId, traceId)
		ctx := metadata.NewIncomingContext(c.Request.Context(), md)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func newTraceId() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
