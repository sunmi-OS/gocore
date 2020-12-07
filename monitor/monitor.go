package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Monitor struct {
	handlerList    []Handler
	closeWaitGroup *sync.WaitGroup
}

type Handler interface {
	SendTextMsg(content string) error
}

// @desc
// @auth liuguoqiang 2020-12-07
// @param
// @return
func NewMonitor(handlerList ...Handler) *Monitor {
	return &Monitor{
		handlerList:    handlerList,
		closeWaitGroup: &sync.WaitGroup{},
	}
}

// @desc 发送本消息
// @auth liuguoqiang 2020-12-07
// @param
// @return
func (monitor *Monitor) SendTextMsg(content string) error {
	monitor.closeWaitGroup.Add(1)
	go func() {
		defer func() {
			monitor.closeWaitGroup.Done()
			if r := recover(); r != nil {
				fmt.Printf("%#v\n", r)
			}
		}()
		for k1 := range monitor.handlerList {
			err := monitor.handlerList[k1].SendTextMsg(content)
			if err != nil {
				fmt.Printf("%#v\n", err)
			}
		}
	}()
	return nil
}

// @desc 程序退出前阻塞直到将数据发送出去,或者超时
// @auth liuguoqiang 2020-12-07
// @param
// @return
func (monitor *Monitor) Close(timeout int64) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		monitor.closeWaitGroup.Wait()
		cancel()
	}(ctx)

	select {
	case <-ctx.Done():
		fmt.Println("monitor safe closed")
		return
	case <-time.After(time.Second * time.Duration(timeout)):
		fmt.Println("monitor timeout!!!")
		return
	}
}
