package otel

import (
	"github.com/gin-gonic/gin"
	otelClient "github.com/sunmi-OS/gocore/v2/lib/tracing/client/otel"
	"github.com/sunmi-OS/gocore/v2/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// ZipkinOtel zipkin+opentelemetry
// serviceName service name
// endPointUrl 链路日志上报地址
// sampleRatio according to the rules when the parent trace has no `SampledFlag`, >= 1 will always sample. < 0 are treated as zero
func ZipkinOtel(serviceName string, endPointUrl string, sampleRatio float64) gin.HandlerFunc {
	appName := serviceName + "#" + utils.GetRunTime()
	_, err := otelClient.InitZipkinTracer(appName, endPointUrl, sampleRatio)
	if err != nil {
		panic(err)
	}
	return otelgin.Middleware(serviceName)
}
