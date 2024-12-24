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
	glog.InfoC(s.c, "[%s %d] [%s] time consume: %s\n", s.f, s.l, s.n, time.Since(s.t).String())
}

// New returns a span to log the time consume.
//
// example:
//
//	span := span.New(ctx)
//	defer span.Finish()
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
