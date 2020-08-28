package xlog

import (
	"fmt"
	"log"
	"os"
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
	e.logger = log.New(os.Stderr, "[ERROR] >> ", log.Lmsgprefix|log.Lshortfile|log.Lmicroseconds|log.Ldate)
}

//func stack() (bs []byte) {
//	var buf [2 << 10]byte
//	runtime.Stack(buf[:], false)
//	return buf[:]
//}
