package span

import (
	"context"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
)

type Span struct {
	c context.Context
	p string
	f string
	t time.Time
}

func (s *Span) Finish() {
	glog.InfoC(s.c, "[%s] [%s] time consume: %s", s.p, s.f, time.Since(s.t).String())
}

func New(ctx context.Context, pkg, fn string) *Span {
	return &Span{c: ctx, p: pkg, f: fn, t: time.Now()}
}
