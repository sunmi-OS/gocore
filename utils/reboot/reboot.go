package reboot

//package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
	"golang.org/x/sync/errgroup"
)


func AutoRestart(ctx context.Context, f func() error) error {
	defer func() {
		var r any
		if r = recover(); r != nil {
			buf := make([]byte, 64*1024)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Fprintf(os.Stderr, "reboot: panic recovered: %s\n%s\n", r, buf)
			glog.Error("panic in reboot proc, err: %s, stack: %s", r, buf)

			AutoRestart(ctx, f) // panic, restart
		}
	}()

	for {
		select {
		case <-ctx.Done():
			glog.InfoC(ctx, "reboot AutoRestart ctx done")
			return ctx.Err()
		default:
			err := f()
			if errors.Is(err, context.Canceled) {
				return nil
			}
			// 如果返回异常，先休息会
			if err != nil {
				time.Sleep(time.Second * 3)
			}
		}
	}
	return nil
}


func AsyncFunc(g *errgroup.Group, f func() error) error {
	defer func() {
		var r any
		if r = recover(); r != nil {
			buf := make([]byte, 64*1024)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Fprintf(os.Stderr, "reboot: panic recovered: %s\n%s\n", r, buf)
			glog.FatalF("panic in reboot proc, err: %s, stack: %s", r, buf)

			restartAsyncLoop(g, f) // panic, restart
		}
	}()

	return f()
}

func restartAsyncLoop(g *errgroup.Group, f func() error) {
	f2 := func() error {
		return AsyncFunc(g, f)
	}
	g.Go(f2)
}
