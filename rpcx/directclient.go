package rpcx

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type DirectClient struct {
	conn *grpc.ClientConn
}

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

// @desc 初始化客户端
// @auth liuguoqiang 2020-04-21
// @param
// @return
func NewDirectClientV2(clientConfig *ClientConfig) (*DirectClient, error) {
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
