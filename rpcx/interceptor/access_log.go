package interceptor

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/sunmi-OS/gocore/v2/utils/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func UnaryAccessLog(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})

	start := time.Now()
	requestDate := start.Format(time.RFC3339)

	var res interface{}
	defer func() {
		md, _ := metadata.FromIncomingContext(ctx)
		ua, clientForwardedIp, traceId := extractFromMD(md)
		clientIp := getPeerAddr(ctx)

		var accessLog = struct {
			Metadata              interface{} `json:"metadata"`
			RequestDate           interface{} `json:"request_date"`
			ProcessTime           interface{} `json:"process_time"`
			ClientIp              interface{} `json:"client_ip"`
			ClientForwardedIp     interface{} `json:"client_forwarded_ip"`
			TraceId               interface{} `json:"traceid"`
			Ua                    interface{} `json:"ua"`
			RequestMethod         interface{} `json:"request_method"`
			RequestUrl            interface{} `json:"request_url"`
			RequestParams         interface{} `json:"request_params"`
			ResponseStatusCode    interface{} `json:"response_status_code"`
			ResponseStatusMessage interface{} `json:"response_status_message"`
			ResponseBody          interface{} `json:"response_body"`
		}{
			md,
			requestDate,
			time.Since(start).Milliseconds(),
			clientIp,
			clientForwardedIp,
			traceId,
			ua,
			"gRPC Unary",
			info.FullMethod,
			req,
			int64(status.Code(err)),
			status.Code(err),
			res,
		}
		d := file.GetPath()
		logFile := d + "/access.log"
		var f *os.File
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		if file.CheckFile(logFile) {
			f, _ = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		} else {
			f, _ = os.Create(logFile)
		}
		b, _ := jsoniter.Marshal(accessLog)
		//log.SetOutput(f)
		_, _ = io.WriteString(f, string(b)+"\r\n")
	}()

	res, err = handler(ctx, req)
	return res, err
}

func StreamAccessLog(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})

	start := time.Now()
	requestDate := start.Format(time.RFC3339)
	ctx := stream.Context()

	defer func() {
		md, _ := metadata.FromIncomingContext(ctx)
		ua, clientForwardedIp, traceId := extractFromMD(md)
		clientIp := getPeerAddr(ctx)

		var accessLog = struct {
			Metadata              interface{} `json:"metadata"`
			RequestDate           interface{} `json:"request_date"`
			ProcessTime           interface{} `json:"process_time"`
			ClientIp              interface{} `json:"client_ip"`
			ClientForwardedIp     interface{} `json:"client_forwarded_ip"`
			TraceId               interface{} `json:"traceid"`
			Ua                    interface{} `json:"ua"`
			RequestMethod         interface{} `json:"request_method"`
			RequestUrl            interface{} `json:"request_url"`
			ResponseStatusCode    interface{} `json:"response_status_code"`
			ResponseStatusMessage interface{} `json:"response_status_message"`
		}{
			md,
			requestDate,
			time.Since(start).Milliseconds(),
			clientIp,
			clientForwardedIp,
			traceId,
			ua,
			"gRPC Stream",
			info.FullMethod,
			int64(status.Code(err)),
			status.Code(err),
		}
		b, _ := jsoniter.Marshal(accessLog)
		d := file.GetPath()
		logFile := d + "/access.log"
		var f *os.File
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		if file.CheckFile(logFile) {
			f, _ = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		} else {
			f, _ = os.Create(logFile)
		}
		_, _ = io.WriteString(f, string(b)+"\r\n")
	}()

	err = handler(srv, stream)
	return err
}

func extractFromMD(md metadata.MD) (ua string, ip string, traceId string) {
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

	return ua, ip, traceId
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
