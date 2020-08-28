package xlog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type InfoLogger struct {
	logger *log.Logger
	once   sync.Once
}

func (i *InfoLogger) logOut(format *string, v ...interface{}) {
	i.once.Do(func() {
		i.new()
	})
	if format != nil {
		i.logger.Output(3, fmt.Sprintf(*format, v...))
		return
	}
	i.logger.Output(3, fmt.Sprintln(v...))
}

func (i *InfoLogger) new() {
	i.logger = log.New(os.Stdout, "[INFO] >> ", log.Lmsgprefix|log.Lshortfile|log.Lmicroseconds|log.Ldate)
}
