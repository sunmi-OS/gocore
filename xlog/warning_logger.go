package xlog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type WarnLogger struct {
	logger *log.Logger
	once   sync.Once
}

func (i *WarnLogger) logOut(format *string, v ...interface{}) {
	i.once.Do(func() {
		i.new()
	})
	if format != nil {
		i.logger.Output(3, fmt.Sprintf(*format, v...))
		return
	}
	i.logger.Output(3, fmt.Sprintln(v...))
}

func (i *WarnLogger) new() {
	i.logger = log.New(os.Stderr, "[WARN] >> ", log.Lmsgprefix|log.Lshortfile|log.Lmicroseconds|log.Ldate)
}
