package gopool

import (
	"github.com/panjf2000/ants/v2"
)

var pool *ants.Pool

// NewPool 初始化goroutine pool
func NewPool(size int, opts ...ants.Option) {
	if size <= 0 {
		size = -1
	}
	p, err := ants.NewPool(size, opts...)
	if err != nil {
		panic(err)
	}
	pool = p
}

// GetPool 获取 pool,如果没有初始化，默认创建容量50的pool
func GetPool() *ants.Pool {
	if pool == nil {
		p, err := ants.NewPool(50)
		if err != nil {
			return nil
		}
		pool = p
	}
	return pool
}

// ClosePool 释放pool
func ClosePool() {
	if pool != nil {
		pool.Release()
	}
}
