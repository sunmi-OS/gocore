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
	defaultCfg := defaultClientConfig(name)
	if gc.cfg == nil {
		gc.cfg = defaultCfg
	} else {
		if cfg.Timeout == "" {
			cfg.Timeout = defaultCfg.Timeout
		}
		if cfg.MaxAttempts == "" {
			cfg.MaxAttempts = defaultCfg.MaxAttempts
		}
		if cfg.InitialBackoff == "" {
			cfg.InitialBackoff = defaultCfg.InitialBackoff
		}
		if cfg.MaxBackoff == "" {
			cfg.MaxBackoff = defaultCfg.MaxBackoff
		}
		if cfg.BackoffMultiplier == "" {
			cfg.BackoffMultiplier = defaultCfg.BackoffMultiplier
		}
		if cfg.WaitForReady == "" {
			cfg.WaitForReady = defaultCfg.WaitForReady
		}
		if cfg.MaxTokens == "" {
			cfg.MaxTokens = defaultCfg.MaxTokens
		}
		if cfg.TokenRatio == "" {
			cfg.TokenRatio = defaultCfg.TokenRatio
		}
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
