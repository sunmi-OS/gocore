package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
)

// A ErrGroup is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid, has no limit on the number of active goroutines,
// and does not cancel on error. use WithContext instead.
type ErrGroup interface {
	// Go calls the given function in a new goroutine.
	//
	// The first call to return a non-nil error cancels the group; its error will be
	// returned by Wait.
	Go(fn func(ctx context.Context) error)

	// GOMAXPROCS set max goroutine to work.
	GOMAXPROCS(n int)

	// Wait blocks until all function calls from the Go method have returned, then
	// returns the first non-nil error (if any) from them.
	Wait() error
}

type group struct {
	wg sync.WaitGroup

	err  error
	once sync.Once

	remain int

	ch    chan func(ctx context.Context) error
	cache []func(ctx context.Context) error

	ctx    context.Context
	cancel context.CancelCauseFunc
}

// WithContext returns a new group with a canceled Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs first.
func WithContext(ctx context.Context) ErrGroup {
	ctx, cancel := context.WithCancelCause(ctx)
	return &group{ctx: ctx, cancel: cancel}
}

func (g *group) GOMAXPROCS(n int) {
	if n <= 0 {
		return
	}
	g.remain = n
	g.ch = make(chan func(context.Context) error)
}

func (g *group) Go(fn func(ctx context.Context) error) {
	g.wg.Add(1)

	if g.ch == nil {
		go g.do(fn)
		return
	}

	select {
	case g.ch <- fn:
	default:
		if g.remain > 0 {
			g.spawn()
		}
		select {
		case g.ch <- fn:
		default:
			g.cache = append(g.cache, fn)
		}
	}
}

func (g *group) Wait() error {
	defer func() {
		select {
		case <-g.ctx.Done():
		default:
			g.cancel(nil)
		}
		if g.ch != nil {
			close(g.ch) // let all receiver exit
		}
	}()

	if g.ch != nil {
		for _, fn := range g.cache {
			g.ch <- fn
		}
	}
	g.wg.Wait()

	return g.err
}

func (g *group) spawn() {
	go func() {
		for fn := range g.ch {
			g.do(fn)
		}
	}()
	g.remain--
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
		err = g.ctx.Err()
	default:
		err = fn(g.ctx)
	}
}
