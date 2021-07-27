package xlog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type DebugLogger struct {
	logger *log.Logger
	once   sync.Once
}

func (i *DebugLogger) logOut(format *string, v ...interface{}) {
	i.once.Do(func() {
		i.new()
	})
	if format != nil {
		i.logger.Output(3, fmt.Sprintf(*format, v...))
		return
	}
	i.logger.Output(3, fmt.Sprintln(v...))
}

func (i *DebugLogger) new() {
	version, _ := strconv.Atoi(strings.Split(runtime.Version(), ".")[1])
	if version >= 14 {
		i.logger = log.New(os.Stdout, "[DEBUG] >> ", 64|log.Lshortfile|log.Ldate|log.Lmicroseconds)
		return
	}
	i.logger = log.New(os.Stdout, "[DEBUG] ", log.Lshortfile|log.Ldate|log.Lmicroseconds)
}
