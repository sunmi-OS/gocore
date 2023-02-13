package interceptor

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"gopkg.in/natefinch/lumberjack.v2"
)

func UnaryAccessLog() grpc.UnaryServerInterceptor {
	logger := initZap("./log/access.log")
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})
		start := time.Now()
		requestDate := start.Format(time.RFC3339)
		var res interface{}
		defer func() {
			md, _ := metadata.FromIncomingContext(ctx)
			_, clientForwardedIp, _, host := extractFromMD(md)
			clientIp := getPeerAddr(ctx)
			fields := []zapcore.Field{
				zap.String("r_time", requestDate),
				zap.Int64("cost", time.Since(start).Milliseconds()),
				zap.String("c_ip", clientIp),
				zap.String("c_f_ip", clientForwardedIp),
				zap.String("schema", "gRPC"),
				zap.String("r_host", host),
				zap.String("r_method", "gRPC/Unary"),
				zap.String("r_q_path", info.FullMethod),
				zap.String("r_path", info.FullMethod),
				zap.Any("r_header", md),
				zap.Any("r_body", req),
				zap.Int("s_status", int(status.Code(err))),
				zap.Any("s_code", int(status.Code(err))),
				zap.Any("s_msg", status.Code(err)),
				zap.Any("s_body", res),
			}
			logger.Info("accesslog", fields...)
		}()
		res, err = handler(ctx, req)
		return res, err
	}
}

func StreamAccessLog() grpc.StreamServerInterceptor {
	logger := initZap("./log/access.log")
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})
		start := time.Now()
		requestDate := start.Format(time.RFC3339)
		ctx := stream.Context()
		defer func() {
			md, _ := metadata.FromIncomingContext(ctx)
			_, clientForwardedIp, _, host := extractFromMD(md)
			fmt.Println(host)
			clientIp := getPeerAddr(ctx)
			fields := []zapcore.Field{
				zap.String("r_time", requestDate),
				zap.Int64("cost", time.Since(start).Milliseconds()),
				zap.String("c_ip", clientIp),
				zap.String("c_f_ip", clientForwardedIp),
				zap.String("schema", "gRPC"),
				zap.String("r_host", host),
				zap.String("r_method", "gRPC/Stream"),
				zap.String("r_q_path", info.FullMethod),
				zap.String("r_path", info.FullMethod),
				zap.Any("r_header", md),
				zap.Any("r_body", ""),
				zap.Int("s_status", int(status.Code(err))),
				zap.Any("s_code", int(status.Code(err))),
				zap.Any("s_msg", status.Code(err)),
				zap.Any("s_body", ""),
			}
			logger.Info("accesslog", fields...)
		}()
		err = handler(srv, stream)
		return err
	}
}

func extractFromMD(md metadata.MD) (ua string, ip string, traceId, host string) {
	if v, ok := md["x-forwarded-user-agent"]; ok {
		ua = fmt.Sprintf("%v", v)
	} else {
		ua = fmt.Sprintf("%v", md["user-agent"])
	}

	if v, ok := md["x-forwarded-for"]; ok && len(v) > 0 {
		ips := strings.Split(v[0], ",")
		ip = ips[0]
	}

	if v, ok := md["x-b3-traceid"]; ok && len(v) > 0 {
		traceId = fmt.Sprintf("%v", v)
	}

	if v, ok := md[":authority"]; ok && len(v) > 0 {
		host = fmt.Sprintf("%v", v[0])
	}

	return ua, ip, traceId, host
}

func getPeerAddr(ctx context.Context) string {
	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
		} else {
			addr = pr.Addr.String()
		}
	}
	return addr
}

// init zap
// fileName log path ./access_log
func initZap(fileName string) *zap.Logger {
	config := zapcore.EncoderConfig{
		MessageKey: "",
		LevelKey:   "",
		TimeKey:    "",
		CallerKey:  "",
	}
	// io.Writer 使用 lumberjack
	infoWriter := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1024, //最大体积，单位M，超过则切割
		MaxBackups: 5,    //最大文件保留数，超过则删除最老的日志文件
		MaxAge:     30,   //最长保存时间30天
		Compress:   true, //是否压缩
	}
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), zap.InfoLevel), //将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.InfoLevel))
}
