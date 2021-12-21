# GLog日志库
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
- glog允许同时注册多个输出端
```go
// SetLogger设置日志打印实例,选择输出到文件,终端,阿里云日志等
func SetLogger(name string, logger logx.GLog) {
	Logger.Store(name, logger)
}
```
- glog 默认使用zap作为底层输出端
```go
//  默认加入zap组件
func init() {
	Logger.Store("zap", &zap.Zap{})
}
```
- 使用示例
```
glog.Debug("zap debug")
glog.Warn("zap warn")
glog.Error("zap error", "呵呵")
```




