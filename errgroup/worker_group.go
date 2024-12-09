package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/sunmi-OS/gocore/v2/worker"
)

type workerGroup struct {
	wg   sync.WaitGroup
	pool worker.Pool

	err  error
	once sync.Once

	ctx    context.Context
	cancel context.CancelCauseFunc
}

// WithWorkerContext returns a new group with a canceled Context derived from ctx which use worker pool.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs first.
func WithWorkerContext(ctx context.Context, pool worker.Pool) Group {
	ctx, cancel := context.WithCancelCause(ctx)

	return &workerGroup{
		pool:   pool,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (g *workerGroup) Go(fn func(ctx context.Context) error) {
	g.wg.Add(1)
	g.do(fn)
}

func (g *workerGroup) GOMAXPROCS(_ int) {}

func (g *workerGroup) Wait() error {
	defer g.cancel(g.err)
	g.wg.Wait()

	return g.err
}

func (g *workerGroup) do(fn func(ctx context.Context) error) {
	g.pool.Go(g.ctx, func(ctx context.Context) {
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
		err = fn(ctx)
	})
}
