package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/sunmi-OS/gocore/v2/worker"
)

type poolgroup struct {
	wg   sync.WaitGroup
	pool worker.Pool

	err  error
	once sync.Once

	ctx    context.Context
	cancel context.CancelCauseFunc
}

// WithPoolContext returns a new group with a canceled Context derived from ctx which use work pool.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs first.
func WithPoolContext(ctx context.Context, pool worker.Pool) Group {
	ctx, cancel := context.WithCancelCause(ctx)

	return &poolgroup{
		pool:   pool,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (pg *poolgroup) Go(fn func(ctx context.Context) error) {
	pg.wg.Add(1)
	pg.do(fn)
}

func (pg *poolgroup) GOMAXPROCS(_ int) {}

func (pg *poolgroup) Wait() error {
	defer pg.cancel(pg.err)
	pg.wg.Wait()

	return pg.err
}

func (pg *poolgroup) do(fn func(ctx context.Context) error) {
	pg.pool.Go(pg.ctx, func(ctx context.Context) {
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("errgroup panic recovered: %+v\n%s", r, string(debug.Stack()))
			}
			if err != nil {
				pg.once.Do(func() {
					pg.err = err
					pg.cancel(err)
				})
			}
			pg.wg.Done()
		}()
		err = fn(ctx)
	})
}
