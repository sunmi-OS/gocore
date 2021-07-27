package close

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

type (
	Close struct {
		Name     string
		Priority int
		Func     func()
	}
	closes []Close
)

var closeHandler closes

func (c closes) Len() int           { return len(c) }
func (c closes) Less(i, j int) bool { return c[i].Priority < c[j].Priority }
func (c closes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func AddShutdown(c ...Close) {
	closeHandler = append(closeHandler, c...)
}

func SignalClose() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	sig := <-c
	sort.Sort(closeHandler)
	if len(closeHandler) > 0 {
		for _, f := range closeHandler {
			fmt.Printf("Close %s ...\n", f.Name)
			f.Func()
		}
	}
	fmt.Printf("Got %s signal. Aborting...\n", sig)
	os.Exit(0)
}
