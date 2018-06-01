package main

import (
	"github.com/jolestar/go-commons-pool"
	"fmt"
	"time"
)

var pCommonPool *pool.ObjectPool

type PoolTest struct{}

func (this *PoolTest) Test() string {
	return "PoolTest"
}

func init() {
	// 初始化连接池配置项
	PoolConfig := pool.NewDefaultPoolConfig()
	// 连接池最大容量设置
	PoolConfig.MaxTotal = 1000
	WithAbandonedConfig := pool.NewDefaultAbandonedConfig()
	// 注册连接池初始化链接方式
	pCommonPool = pool.NewObjectPoolWithAbandonedConfig(pool.NewPooledObjectFactorySimple(
		func() (interface{}, error) {
			return Link()
		}), PoolConfig, WithAbandonedConfig)
}

// 初始化链接类
func Link() (*PoolTest, error) {
	fmt.Println("初始化PoolTest类!!!")
	return &PoolTest{}, nil
}

func main() {

	//----------------------------------第一次使用将会调用初始化方法---------------------------------
	fmt.Println("第一次使用将会调用初始化方法")
	Test()

	//----------------------------------第二次使用将会复用初始化好的对象---------------------------------
	fmt.Println("第二次使用将会复用初始化好的实例")
	Test()

	//----------------------------------连续多次并发调用当连接池不够用的会扩充连接池---------------------------
	fmt.Println("连续多次并发调用当连接池不够用的会扩充连接池")
	go Test()
	go Test()
	go Test()
	go Test()
	go Test()

	time.Sleep(1 * time.Second)
}

func Test() {
	var client *PoolTest
	// 从连接池中获取一个实例
	obj, _ := pCommonPool.BorrowObject()
	// 转换为对应实体
	if obj != nil {
		client = obj.(*PoolTest)
	}
	// 调用需要的方法
	fmt.Println(client.Test())
	// 交还连接池
	pCommonPool.ReturnObject(client)
}
