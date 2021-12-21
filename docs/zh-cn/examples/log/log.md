# GLog日志库
- glog 默认使用zap作为底层输出端,
- glog允许同时注册多个输出端
- 只要实现logx中的GLog interface{}即可注册到glog的输出端
```go
type GLog interface {
	Info(args ...interface{})
	InfoF(format string, args ...interface{})
	Debug(args ...interface{})
	DebugF(format string, args ...interface{})
	Warn(args ...interface{})
	WarnF(format string, args ...interface{})
	Error(args ...interface{})
	ErrorF(format string, args ...interface{})
}
```



