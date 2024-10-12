package reboot

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
	"golang.org/x/sync/errgroup"
)

func startBooter(g *errgroup.Group, f func() error) {

	f2 := func() error {
		return AsyncFunc(g, f)
	}
	g.Go(f2)
}

func TestBootWithErrgroup(t *testing.T) {
	var g errgroup.Group

	startBooter(&g, bizFunc)

	if err := g.Wait(); err != nil {
		fmt.Println("err:", err)
	}

	time.Sleep(20 * time.Second)

}

func bizFunc() error {
	// manual make panic for test
	for i := 0; i < 10; i++ {
		glog.InfoF("running in goroutine:%d", i)
		time.Sleep(1 * time.Second)
		if i == 5 {
			panic("just boom")
		}
	}
	return nil
}

// use method
func TestBootAutoRestart(t *testing.T) {
	ctx, canel := context.WithCancel(context.Background())
	go AutoRestart(ctx, bizFunc)

	time.Sleep(20 * time.Second)
	canel()
	time.Sleep(20 * time.Second)
}
