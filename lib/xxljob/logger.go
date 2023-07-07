package xxljob

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type LogLevel string

const (
	//DebugLevel = "debug"
	InfoLevel = "info"
	//WarnLevel  = "warn"
	ErrorLevel = "error"
)

const (
	errorLevel int = iota + 1
	warnLevel
	infoLevel
	debugLevel
)

type logger struct {
	logger    *log.Logger
	once      sync.Once
	level     LogLevel
	levelInt  int
	CallDepth int
}

func newLogger(level LogLevel) *logger {
	l := &logger{
		level:     level,
		CallDepth: 3, // default call depth
	}
	l.configLogger()
	return l
}

func (l *logger) configLogger() {
	switch l.level {
	//case DebugLevel:
	//	l.levelInt = debugLevel
	//	l.logger = log.New(os.Stdout, "[DEBUG] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	case InfoLevel:
		l.levelInt = infoLevel
		l.logger = log.New(os.Stdout, "[INFO] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	//case WarnLevel:
	//	l.levelInt = warnLevel
	//	l.logger = log.New(os.Stdout, "[WARN] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	case ErrorLevel:
		l.levelInt = errorLevel
		l.logger = log.New(os.Stdout, "[ERROR] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	default:
		l.levelInt = infoLevel
		l.logger = log.New(os.Stdout, "[INFO] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	}
}

func (l *logger) Info(format string, args ...interface{}) {
	if l.level >= InfoLevel {
		if &format != nil {
			_ = l.logger.Output(l.CallDepth, fmt.Sprintf(format, args...))
			return
		}
		_ = l.logger.Output(l.CallDepth, fmt.Sprintln(args...))
	}
}

func (l *logger) Error(format string, args ...interface{}) {
	if l.level >= ErrorLevel {
		if &format != nil {
			_ = l.logger.Output(l.CallDepth, fmt.Sprintf(format, args...))
			return
		}
		_ = l.logger.Output(l.CallDepth, fmt.Sprintln(args...))
	}
}
