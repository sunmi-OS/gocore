package errgroup

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnlimitedRunsConcurrently(t *testing.T) {
	const tasks = 8

	eg := WithContext(context.Background(), 0)

	// 无限流模式：所有任务必须同时在跑才能返回，否则超时
	var barrier sync.WaitGroup
	barrier.Add(tasks)
	for range tasks {
		eg.Go(func(context.Context) error {
			barrier.Done()
			barrier.Wait()
			return nil
		})
	}

	finished := make(chan error, 1)
	go func() { finished <- eg.Wait() }()

	select {
	case err := <-finished:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("tasks did not run concurrently")
	}
}

func TestFirstErrorCancelsGroup(t *testing.T) {
	errBoom := errors.New("boom")

	eg := WithContext(context.Background(), 2)

	started := make(chan struct{})
	cause := make(chan error, 1)

	eg.Go(func(ctx context.Context) error {
		close(started)
		<-ctx.Done() // 等待另一个任务的错误取消整个 group
		cause <- context.Cause(ctx)
		return nil
	})
	eg.Go(func(context.Context) error {
		<-started // 确保另一个任务已开始执行
		return errBoom
	})

	assert.ErrorIs(t, eg.Wait(), errBoom)

	select {
	case err := <-cause:
		assert.ErrorIs(t, err, errBoom)
	case <-time.After(time.Second):
		t.Fatal("running task did not observe cancellation")
	}
}

func TestLimitTasksRunBeforeWait(t *testing.T) {
	eg := WithContext(context.Background(), 1)
	release := make(chan struct{})
	started := make(chan struct{})

	eg.Go(func(context.Context) error {
		<-release
		return nil
	})
	eg.Go(func(context.Context) error {
		close(started)
		return nil
	})

	close(release)
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("queued task did not run before Wait")
	}

	assert.NoError(t, eg.Wait())
}

func TestConcurrentGoWithLimit(t *testing.T) {
	const (
		limit = 4
		tasks = 100
	)

	eg := WithContext(context.Background(), limit)
	start := make(chan struct{})

	var submit sync.WaitGroup
	submit.Add(tasks)

	var active int32
	var maxActive int32
	var completed int32

	for range tasks {
		go func() {
			defer submit.Done()
			<-start

			eg.Go(func(context.Context) error {
				current := atomic.AddInt32(&active, 1)
				for {
					max := atomic.LoadInt32(&maxActive)
					if current <= max || atomic.CompareAndSwapInt32(&maxActive, max, current) {
						break
					}
				}
				time.Sleep(time.Millisecond)
				atomic.AddInt32(&active, -1)
				atomic.AddInt32(&completed, 1)
				return nil
			})
		}()
	}

	close(start)
	submit.Wait()

	assert.NoError(t, eg.Wait())
	assert.Equal(t, int32(tasks), atomic.LoadInt32(&completed))
	assert.LessOrEqual(t, atomic.LoadInt32(&maxActive), int32(limit))
}

func TestLimitReuseIdleWorkers(t *testing.T) {
	g := WithContext(context.Background(), 100).(*group)

	for range 10 {
		done := make(chan struct{})
		g.Go(func(context.Context) error {
			close(done)
			return nil
		})

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("task did not run")
		}
		waitErrGroupIdle(t, g, 1)
	}

	g.mutex.Lock()
	spawned := 100 - g.remain
	g.mutex.Unlock()

	assert.Equal(t, 1, spawned)
	assert.NoError(t, g.Wait())
}

func TestBurstSpawnsEnoughWorkers(t *testing.T) {
	const limit = 4

	g := WithContext(context.Background(), limit).(*group)

	// 先养出一个空闲 worker
	done := make(chan struct{})
	g.Go(func(context.Context) error {
		close(done)
		return nil
	})
	<-done
	waitErrGroupIdle(t, g, 1)

	// 突发提交 4 个任务：全部同时在跑才能返回，
	// 若并发塌缩（唤醒窗口内不扩容），测试会超时
	var running sync.WaitGroup
	running.Add(limit)
	for range limit {
		g.Go(func(context.Context) error {
			running.Done()
			running.Wait()
			return nil
		})
	}

	finished := make(chan error, 1)
	go func() { finished <- g.Wait() }()

	select {
	case err := <-finished:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("burst tasks did not run concurrently")
	}
}

func TestGoAfterWaitPanicsWithLimit(t *testing.T) {
	eg := WithContext(context.Background(), 1)
	assert.NoError(t, eg.Wait())
	assert.Panics(t, func() {
		eg.Go(func(context.Context) error { return nil })
	})
}

func TestRecover(t *testing.T) {
	eg := WithContext(context.Background(), 1)
	eg.Go(func(context.Context) (err error) {
		panic("oh my god!")
	})

	err := eg.Wait()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "panic recovered"), "error should contain panic info: %v", err)
	assert.True(t, strings.Contains(err.Error(), "oh my god!"), "error should contain panic value: %v", err)
}

func TestZeroGroup(t *testing.T) {
	err1 := errors.New("errgroup_test: 1")
	err2 := errors.New("errgroup_test: 2")

	cases := []struct {
		errs []error
	}{
		{errs: []error{}},
		{errs: []error{nil}},
		{errs: []error{err1}},
		{errs: []error{err1, nil}},
		{errs: []error{err1, nil, err2}},
	}

	for _, tc := range cases {
		eg := WithContext(context.Background(), 0)

		var firstErr error
		for i, err := range tc.errs {
			err := err
			eg.Go(func(context.Context) error { return err })

			if firstErr == nil && err != nil {
				firstErr = err
			}

			if gErr := eg.Wait(); gErr != firstErr {
				t.Errorf("after g.Go(func() error { return err }) for err in %v\n"+
					"g.Wait() = %v; want %v", tc.errs[:i+1], err, firstErr)
			}
		}
	}
}

func waitErrGroupIdle(t *testing.T, g *group, n int) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		g.mutex.Lock()
		idle := g.idle
		g.mutex.Unlock()

		if idle >= n {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("timed out waiting for %d idle worker(s)", n)
}
