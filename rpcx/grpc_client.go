package rpcx

import (
	"google.golang.org/grpc"
)

type GrpcClient struct {
	Name       string
	addr       string
	ClientConn *grpc.ClientConn
	cfg        *GrpcClientConfig
}

// NewGrpcClient new grpc client
func NewGrpcClient(name, addr string, cfg *GrpcClientConfig) (gc *GrpcClient, err error) {

	gc = &GrpcClient{cfg: cfg, addr: addr}
	if gc.cfg == nil {
		gc.cfg = defaultClientConfig(name)
	}

	options := gc.cfg.buildDialOptions()
	conn, err := grpc.Dial(gc.addr, options...)
	if err != nil {
		return nil, err
	}
	gc.ClientConn = conn
	return gc, nil
}

func (c *GrpcClient) Close() {
	if c.ClientConn != nil {
		c.ClientConn.Close()
	}
}
