package otel

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	gormSpanKey        = "__gorm_span"
	callBackBeforeName = "opentelemetry:before"
	callBackAfterName  = "opentelemetry:after"
)

func before(db *gorm.DB) {
	// 先从父级spans生成子span
	_, span := otel.Tracer("").Start(db.Statement.Context, db.Statement.Table)
	// 利用db实例去传递span
	db.InstanceSet(gormSpanKey, span)
}

func after(db *gorm.DB) {
	// 从GORM的DB实例中取出span
	_span, isExist := db.InstanceGet(gormSpanKey)
	if !isExist {
		return
	}
	// 断言进行类型转换
	span, ok := _span.(trace.Span)
	if !ok {
		return
	}
	defer span.End()
	// span.SetTag("db.statement", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...))
	attr := attribute.KeyValue{
		Key:   "db.statement",
		Value: attribute.StringValue(db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)),
	}
	span.SetAttributes(attr)
}

type OpentelemetryPlugin struct{}

func (op *OpentelemetryPlugin) Name() string {
	return "opentelemetryPlugin"
}

func (op *OpentelemetryPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)
	// 结束后
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

var _ gorm.Plugin = &OpentelemetryPlugin{}

// StartSpanWithCtx 生成上下文span
// skip The argument skip is the number of stack frames to ascend, with 0 identifying the caller of Caller
func StartSpanWithCtx(ctx context.Context, db *gorm.DB, skip int) (trace.Span, *gorm.DB) {
	_ = db.Use(&OpentelemetryPlugin{})
	pc, _, _, _ := runtime.Caller(skip)
	spanName := "/" + runtime.FuncForPC(pc).Name()
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	db = db.WithContext(ctx)
	return span, db
}
