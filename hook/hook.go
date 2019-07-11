package hook

import (
	"container/list"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Hook struct {
	task *list.List
	lock sync.Mutex
}

var hookHandler = Hook{}

func AddShutdownHook(runnables ...func() int)  {
	defer hookHandler.lock.Unlock()
	hookHandler.lock.Lock()
	if hookHandler.task == nil {
		hookHandler.task = new(list.List)
	}
	for _, v := range runnables {
		hookHandler.task.PushBack(v)
	}
	if len(runnables) > 0 && len(runnables) == hookHandler.task.Len() {
		go hookHandler.listenShutdownSignal() // only start once
	}
}


func (h *Hook) listenShutdownSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-c:
		replicasTask := h.taskReplicas()
		if replicasTask != nil && len(replicasTask) > 0 {
			exitCode := replicasTask[0]()
			for i := 1; i < len(replicasTask); i++ {
				replicasTask[i]()
			}
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			os.Exit(exitCode)
		}
	}
}

func (h *Hook) taskReplicas() []func() int {
	defer h.lock.Unlock()
	h.lock.Lock()
	replicas := make([]func() int, h.task.Len())
	i := 0
	for e := h.task.Front(); e != nil; e=e.Next() {
		replicas[i] = e.Value.(func() int)
		i++
	}
	return replicas
}
