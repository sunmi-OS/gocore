package gopool

import (
	"testing"

	"github.com/panjf2000/ants/v2"
)

func TestNewPool(t *testing.T) {
	NewPool(10)
	if pool == nil {
		t.Error("NewPool error")
	}
}

func TestGetPool(t *testing.T) {
	p := GetPool()
	if p == nil {
		t.Error("GetPool error")
	}
	if p.Cap() != 50 {
		t.Error("GetPool error")
	}
	pool, _ = ants.NewPool(10)
	if GetPool().Cap() != pool.Cap() {
		t.Error("GetPool error")
	}
}

func TestClosePool(t *testing.T) {
	pool, _ = ants.NewPool(10)
	ClosePool()
	if !pool.IsClosed() {
		t.Error("ClosePool error")
	}
}
