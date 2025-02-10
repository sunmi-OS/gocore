package worker

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	RunTimes   = 1e6
	PoolCap    = 5e4
	BenchParam = 10
)

func demoFunc(ctx context.Context) {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

func BenchmarkGoroutines(b *testing.B) {
	runtime.GOMAXPROCS(1)
	ctx := context.TODO()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			go func() {
				demoFunc(ctx)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkChannel(b *testing.B) {
	runtime.GOMAXPROCS(1)
	ctx := context.TODO()

	var wg sync.WaitGroup
	sema := make(chan struct{}, PoolCap)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoFunc(ctx)
				<-sema
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkErrGroup(b *testing.B) {
	runtime.GOMAXPROCS(1)
	ctx := context.TODO()

	var wg sync.WaitGroup
	var eg errgroup.Group
	eg.SetLimit(PoolCap)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			eg.Go(func() error {
				demoFunc(ctx)
				wg.Done()
				return nil
			})
		}
		wg.Wait()
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	runtime.GOMAXPROCS(1)
	ctx := context.TODO()

	p := NewPool(PoolCap)
	defer p.Close()

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			p.Go(ctx, func(ctx context.Context) {
				demoFunc(ctx)
				wg.Done()
			})
		}
		wg.Wait()
	}
}

func BenchmarkGoroutinesThroughput(b *testing.B) {
	runtime.GOMAXPROCS(1)
	ctx := context.TODO()

	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			go demoFunc(ctx)
		}
	}
}

func BenchmarkSemaphoreThroughput(b *testing.B) {
	runtime.GOMAXPROCS(1)
	ctx := context.TODO()

	sema := make(chan struct{}, PoolCap)
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoFunc(ctx)
				<-sema
			}()
		}
	}
}

func BenchmarkWorkerPoolThroughput(b *testing.B) {
	runtime.GOMAXPROCS(1)
	ctx := context.TODO()

	p := NewPool(PoolCap)
	defer p.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			p.Go(ctx, demoFunc)
		}
	}
}
