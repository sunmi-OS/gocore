package rpcx

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// Deprecated
type DirectClient struct {
	conn *grpc.ClientConn
	cfg  *ClientConfig
}

type GrpcClient struct {
	Name       string
	addr       string
	ClientConn *grpc.ClientConn
	cfg        *GrpcClientConfig
}

// Deprecated
// 推荐使用 NewGrpcClient
// @desc 初始化客户端 timeout:单位毫秒,默认两秒超时
// @auth liuguoqiang 2020-04-21
// @param
// @return
func NewDirectClient(host string, timeout int64, opts ...ClientOption) (*DirectClient, error) {
	options := []ClientOption{}
	if timeout > 0 {
		options = append(options, WithTimeout(time.Duration(timeout)*time.Millisecond))
	}
	options = append(options, opts...)

	ops := buildDialOptions(options...)
	conn, err := grpc.Dial(host, ops...)
	if err != nil {
		return nil, err
	}

	return &DirectClient{
		conn: conn,
	}, nil
}

// Deprecated
// 推荐使用 NewGrpcClient
// @desc 初始化客户端
// @auth liuguoqiang 2020-04-21
// @param
// @return
func NewDirectClientV2(clientConfig *ClientConfig) (*DirectClient, error) {
	if clientConfig == nil || clientConfig.Host == "" {
		return nil, fmt.Errorf("rpc链接地址不能为空")
	}
	options := clientConfig.buildDialOptions()
	conn, err := grpc.Dial(clientConfig.Host, options...)
	if err != nil {
		return nil, err
	}
	return &DirectClient{
		conn: conn,
	}, nil
}

// @desc 返回grpc链接
// @auth liuguoqiang 2020-04-21
// @param
// @return
func (c *DirectClient) Next() (*grpc.ClientConn, bool) {
	state := c.conn.GetState()
	if state == connectivity.Ready {
		return c.conn, true
	} else {
		return c.conn, false
	}
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

// Conn get ready *grpc.ClientConn
func (c *GrpcClient) Conn() (conn *grpc.ClientConn, ok bool) {
	state := c.ClientConn.GetState()
	if state == connectivity.Ready {
		ok = true
	}
	return c.ClientConn, ok
}

func (c *GrpcClient) Close() {
	if c.ClientConn != nil {
		c.ClientConn.Close()
	}
}
