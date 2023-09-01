package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func GetMetaData(ctx context.Context, key string) string {
	if md, b := metadata.FromIncomingContext(ctx); b {
		vals := md.Get(key)
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

func SetMetaData(ctx context.Context, key string, val string) context.Context {
	md, b := metadata.FromIncomingContext(ctx)
	if !b {
		md = metadata.MD{}
	}
	md.Set(key, val)
	return metadata.NewIncomingContext(ctx, md)
}

func SetMetaDataMulti(ctx context.Context, kvs map[string]string) context.Context {
	md, b := metadata.FromIncomingContext(ctx)
	if !b {
		md = metadata.MD{}
	}
	for k, v := range kvs {
		md.Set(k, v)
	}
	return metadata.NewIncomingContext(ctx, md)
}

func GetMetaDataMulti(ctx context.Context, keys []string) map[string]string {
	md, b := metadata.FromIncomingContext(ctx)
	if !b {
		return nil
	}
	res := make(map[string]string)
	for _, k := range keys {
		vals := md.Get(k)
		if len(vals) > 0 {
			res[k] = vals[0]
		} else {
			res[k] = ""
		}
	}
	return res
}
