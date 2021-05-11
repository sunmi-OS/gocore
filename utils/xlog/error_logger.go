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

type ErrorLogger struct {
	logger *log.Logger
	once   sync.Once
}

func (e *ErrorLogger) logOut(format *string, v ...interface{}) {
	e.once.Do(func() {
		e.new()
	})
	if format != nil {
		e.logger.Output(3, fmt.Sprintf(*format, v...))
		//i.logger.Writer().Write(stack())
		return
	}
	e.logger.Output(3, fmt.Sprintln(v...))
	//i.logger.Writer().Write(stack())
}

func (e *ErrorLogger) new() {
	version, _ := strconv.Atoi(strings.Split(runtime.Version(), ".")[1])
	if version >= 14 {
		e.logger = log.New(os.Stdout, "[ERROR] >> ", 64|log.Lshortfile|log.Ldate|log.Lmicroseconds)
		return
	}
	e.logger = log.New(os.Stdout, "[ERROR] ", log.Lshortfile|log.Ldate|log.Lmicroseconds)
}
