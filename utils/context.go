package utils

import (
	"context"
	"time"
)

// WithoutCancel returns a copy of parent that is not canceled when parent is canceled.
// The returned context returns no Deadline or Err, and its Done channel is nil.
func WithoutCancel(parent context.Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return withoutCancelCtx{parent}
}

type withoutCancelCtx struct {
	context.Context
}

func (withoutCancelCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (withoutCancelCtx) Done() <-chan struct{} {
	return nil
}

func (withoutCancelCtx) Err() error {
	return nil
}
