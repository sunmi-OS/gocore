package otel

import (
	"context"
	"runtime"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// InitZipkinTracer ipkin exporter for OpenTelemetry
// sampleRatio according to the rules when the parent trace has no `SampledFlag`, >= 1 will always sample. < 0 are treated as zero
func InitZipkinTracer(appName, endPointUrl string, sampleRatio float64) (*sdktrace.TracerProvider, error) {
	exporter, err := zipkin.New(
		endPointUrl,
	)
	if err != nil {
		return nil, err
	}
	res, err := resource.New(context.Background(),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appName),
		),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		// sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRatio))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))
	return tp, nil
}

// StartSpan Generate span based on context
func StartSpan(ctx context.Context) (context.Context, trace.Span) {
	// skip The argument skip is the number of stack frames to ascend, with 0 identifying the caller of Caller
	pc, _, _, _ := runtime.Caller(1)
	spanName := "/" + runtime.FuncForPC(pc).Name()
	return otel.Tracer("").Start(ctx, spanName)
}

// StartSpanWithName Generate named span based on context
func StartSpanWithName(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, spanName)
}
