package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
)

// A ErrGroup is a collection of goroutines working on subtasks that are part of
// the same overall task. A ErrGroup should not be reused for different tasks.
//
// Use WithContext to create a ErrGroup.
type ErrGroup interface {
	// Go calls the given function in a goroutine.
	// Go never blocks; if the concurrency limit is reached, the function is
	// queued and executed as running goroutines complete.
	//
	// The first call to return a non-nil error cancels the group; its error will be
	// returned by Wait.
	//
	// Go may be called concurrently, but all calls must happen before Wait
	// returns. Calling Go after Wait has returned panics when limit > 0.
	Go(fn func(ctx context.Context) error)

	// Wait blocks until all function calls from the Go method have returned, then
	// returns the first non-nil error (if any) from them.
	//
	// Wait must not be called concurrently with external Go calls that may
	// increase the task count from zero (the same constraint as sync.WaitGroup).
	Wait() error
}

type group struct {
	wg sync.WaitGroup

	err  error
	once sync.Once

	mutex  sync.Mutex
	cond   *sync.Cond // 仅限制并发数时初始化
	tasks  []func(ctx context.Context) error
	remain int
	idle   int
	closed bool

	ctx    context.Context
	cancel context.CancelCauseFunc
}

// WithContext returns a new ErrGroup that is associated with a derived Context.
//
// The returned group's Context is canceled in the following cases:
//   - The first time a goroutine started with Go returns a non-nil error.
//   - Or when Wait is called and returns.
//
// If limit > 0, the group restricts the number of active goroutines
// to at most 'limit'. Additional functions passed to Go will be queued
// and executed only when running goroutines complete.
//
// The derived Context is created with context.WithCancelCause, so the
// cancellation reason is preserved and can be retrieved via context.Cause.
func WithContext(ctx context.Context, limit int) ErrGroup {
	ctx, cancel := context.WithCancelCause(ctx)

	g := &group{
		ctx:    ctx,
		cancel: cancel,
	}
	if limit > 0 {
		g.remain = limit
		g.cond = sync.NewCond(&g.mutex)
	}
	return g
}

func (g *group) Go(fn func(ctx context.Context) error) {
	// 不限制并发数，直接新开协程
	if g.cond == nil {
		g.wg.Add(1)
		go g.do(fn)
		return
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Wait 已返回，协程均已退出，任务不会被执行
	if g.closed {
		panic("errgroup: Go called after Wait")
	}

	g.wg.Add(1)
	g.tasks = append(g.tasks, fn)
	// 积压任务超过空闲协程消费能力且未达上限时，新开一个协程；
	// 否则复用空闲协程
	if len(g.tasks) > g.idle && g.remain > 0 {
		g.remain--
		g.spawn()
	}
	// 唤醒一个空闲协程
	g.cond.Signal()
}

func (g *group) Wait() error {
	g.wg.Wait()
	// 所有任务完成后取消派生context
	g.cancel(nil)

	if g.cond != nil {
		g.mutex.Lock()
		g.closed = true
		g.mutex.Unlock()

		// 通知所有协程退出
		g.cond.Broadcast()
	}

	return g.err
}

func (g *group) spawn() {
	go func() {
		for {
			g.mutex.Lock()

			// 当前无任务且未关闭，进入等待；唤醒后重新检查任务和关闭状态
			for len(g.tasks) == 0 {
				if g.closed {
					g.mutex.Unlock()
					return
				}

				// 进入等待
				g.idle++
				g.cond.Wait()
				g.idle--
			}

			// 取出一个任务
			fn := g.tasks[0]
			g.tasks[0] = nil // 避免引用滞留
			g.tasks = g.tasks[1:]
			g.mutex.Unlock()

			// 执行任务
			g.do(fn)
		}
	}()
}

func (g *group) do(fn func(ctx context.Context) error) {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("errgroup panic recovered: %+v\n%s", r, string(debug.Stack()))
		}
		if err != nil {
			g.once.Do(func() {
				g.err = err
				g.cancel(err)
			})
		}
		g.wg.Done()
	}()

	select {
	case <-g.ctx.Done():
		err = context.Cause(g.ctx)
	default:
		err = fn(g.ctx)
	}
}
