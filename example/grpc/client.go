package proto

import (
	"context"
	"errors"
	"github.com/sunmi-OS/gocore/rpcx"
	"golang-example/pkg/istio"
	"google.golang.org/grpc/metadata"
	"log"
)

var Rpc *rpcx.DirectClient

func Init(host string, timeout int64, opts ...rpcx.ClientOption) {

	var err error
	Rpc, err = rpcx.NewDirectClient(host, timeout, opts...)
	if err != nil {
		log.Fatal("rpc connect fail")
	}
}

func PrintOk(in *Request, trace istio.TraceHeader) (*Response, error) {
	conn, ok := Rpc.Next()
	if !ok || conn == nil {
		return nil, errors.New("rpc connect fail")
	}
	client := NewPrintServiceClient(conn)

	ctx := metadata.NewOutgoingContext(context.Background(), trace.Grpc_MD)
	return client.PrintOK(ctx, in)
}
