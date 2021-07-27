package rate

import (
	"sync"
)

// RateGroup 速率Group，懒加载
type RateGroup struct {
	new func() interface{}
	rgs sync.Map
	sync.RWMutex
}

// NewRateGroup 新建RateGroup
func NewRateGroup(new func() interface{}) (rg *RateGroup) {
	if new == nil {
		panic("RateGroup: can't assign a nil to the new function")
	}
	return &RateGroup{new: new}
}

// Get 获取RateGroup，如果没有则新建
func (r *RateGroup) Get(key string) interface{} {
	rg, ok := r.rgs.Load(key)
	if !ok {
		r.RLock()
		newRg := r.new
		r.RLock()
		rg = newRg()
		r.rgs.Store(key, rg)
	}
	return rg
}
