package xlog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type WarningLogger struct {
	logger *log.Logger
	once   sync.Once
}

func (i *WarningLogger) logOut(format *string, v ...interface{}) {
	i.once.Do(func() {
		i.init()
	})
	if format != nil {
		i.logger.Output(3, fmt.Sprintf(*format, v...))
		return
	}
	i.logger.Output(3, fmt.Sprintln(v...))
}

func (i *WarningLogger) init() {
	i.logger = log.New(os.Stderr, "[WARNING] >> ", log.Lmsgprefix|log.Lshortfile|log.Lmicroseconds|log.Ldate)
}
