package hbase

import (
	"net"
	"sync"

	"github.com/jolestar/go-commons-pool"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/spf13/cast"
	"strings"
)

var onceHbaseClient sync.Once
var HbaseClinet *RpcHbaseClient
var HbasePool *pool.ObjectPool

type RpcHbaseClient struct{}

func NewHbase() error {

	// 验证连通性
	obj, err := LinkHbase()
	if err != nil {
		return err
	}
	// 验证完成关闭通道
	obj.Transport.Close()

	onceHbaseClient.Do(func() {
		PoolConfig := pool.NewDefaultPoolConfig()
		PoolConfig.MaxTotal = cast.ToInt(viper.GetEnvConfig("hbase.PoolSum"))
		WithAbandonedConfig := pool.NewDefaultAbandonedConfig()
		HbasePool = pool.NewObjectPoolWithAbandonedConfig(pool.NewPooledObjectFactorySimple(
			func() (interface{}, error) {
				return LinkHbase()
			}), PoolConfig, WithAbandonedConfig)

		HbaseClinet = &RpcHbaseClient{}
	})

	return nil

}

func LinkHbase() (*THBaseServiceClient, error) {

	// 初始化Thrift建立连接
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err := thrift.NewTSocket(net.JoinHostPort(viper.GetEnvConfig("hbase.Host"), viper.GetEnvConfig("hbase.Port")))

	if err != nil {
		return nil, err
	}

	client := NewTHBaseServiceClientFactory(transport, protocolFactory)

	if err := transport.Open(); err != nil {
		return nil, err
	}

	return client, nil
}

// --------------------------------具体函数实现------------------------------------

func (h RpcHbaseClient) Get(table []byte, tget *TGet) (r *TResult_, err error) {

	var client *THBaseServiceClient

	obj, err := HbasePool.BorrowObject()
	if obj == nil || err != nil {
		return
	}
	client = obj.(*THBaseServiceClient)

	r, err = client.Get(table, tget)
	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "Connection not open") {
			// thrift服务破碎需要重连
			client, err = LinkHbase()
			if err != nil {
				return
			}
			r, err = client.Get(table, tget)
		} else {
			return
		}
	}

	HbasePool.ReturnObject(client)
	return
}

func (h RpcHbaseClient) Exists(table []byte, tget *TGet) (r bool, err error) {

	var client *THBaseServiceClient

	obj, err := HbasePool.BorrowObject()
	if obj == nil || err != nil {
		return
	}
	client = obj.(*THBaseServiceClient)

	r, err = client.Exists(table, tget)
	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "Connection not open") {
			// thrift服务破碎需要重连
			client, err = LinkHbase()
			if err != nil {
				return
			}
			r, err = client.Exists(table, tget)
		} else {
			return
		}
	}

	HbasePool.ReturnObject(client)
	return
}

func (h RpcHbaseClient) Put(table []byte, tput *TPut) (err error) {

	var client *THBaseServiceClient

	obj, err := HbasePool.BorrowObject()
	if obj == nil || err != nil {
		return
	}
	client = obj.(*THBaseServiceClient)

	err = client.Put(table, tput)
	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "Connection not open") {
			// thrift服务破碎需要重连
			client, err = LinkHbase()
			if err != nil {
				return
			}
			err = client.Put(table, tput)
		} else {
			return
		}
	}

	HbasePool.ReturnObject(client)
	return
}




func (h RpcHbaseClient) Delete(table []byte, tdelete *TDelete) (err error) {

	var client *THBaseServiceClient

	obj, err := HbasePool.BorrowObject()
	if obj == nil || err != nil {
		return
	}
	client = obj.(*THBaseServiceClient)

	err = client.DeleteSingle(table, tdelete)
	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "Connection not open") {
			// thrift服务破碎需要重连
			client, err = LinkHbase()
			if err != nil {
				return
			}
			err = client.DeleteSingle(table, tdelete)
		} else {
			return
		}
	}

	HbasePool.ReturnObject(client)
	return
}

