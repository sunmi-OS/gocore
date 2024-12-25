package span

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
)

type Span struct {
	c context.Context
	f string
	l int
	n string
	t time.Time
}

func (s *Span) Finish() {
	glog.InfoC(s.c, "[time consume] function=%s, duration=%s, file=%s:%d", s.n, time.Since(s.t).String(), s.f, s.l)
}

// New returns a span to log the time consume.
//
// example:
//
//	sp := span.New(ctx)
//	defer sp.Finish()
func New(ctx context.Context) *Span {
	sp := &Span{c: ctx, t: time.Now()}
	// Skip level 1 to get the caller function
	pc, file, line, _ := runtime.Caller(1)
	sp.f, sp.l = file, line
	// Get the function details
	if fn := runtime.FuncForPC(pc); fn != nil {
		name := fn.Name()
		sp.n = name[strings.Index(name, ".")+1:]
	}
	return sp
}
