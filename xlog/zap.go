package xlog

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sunmi-OS/gocore/utils"
	"github.com/sunmi-OS/gocore/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	once        sync.Once
	Logger      *zap.Logger
	Sugar       *zap.SugaredLogger
	c           zap.Config
	logFile     *os.File
	logFileDir  string
	logFileName string
	err         error
}

var z = &Logger{}

func Zap() *Logger {
	z.once.Do(func() {
		z.initZap()
	})
	return z
}

func (l *Logger) Info(args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Info(args...)
		return
	}
	infoLog.logOut(nil, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Infof(format, args...)
		return
	}
	infoLog.logOut(&format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Debug(args...)
		return
	}
	debugLog.logOut(nil, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Debugf(format, args...)
		return
	}
	debugLog.logOut(&format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Warn(args...)
		return
	}
	warnLog.logOut(nil, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Warnf(format, args...)
		return
	}
	warnLog.logOut(&format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Error(args...)
		return
	}
	errLog.logOut(nil, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Errorf(format, args...)
		return
	}
	errLog.logOut(&format, args...)
}

func (l *Logger) initZap() {
	// new log file dir
	l.newLogDir()
	// new log file
	l.newLogFile()
	// new zap
	err := l.newZap()
	if err != nil {
		Errorf("l.newZap(),err:%+v", err)
		return
	}
	// updateLogFile Loop
	go l.updateLogFile()
}

func (l *Logger) newZap() (err error) {
	l.c = zap.NewProductionConfig()
	if l.err == nil && l.logFileName != "" {
		l.c.OutputPaths = []string{l.logFileName, "stdout"}
		l.c.ErrorOutputPaths = []string{l.logFileName, "stderr"}
	}
	l.c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	viper.C.SetDefault("system.debug", "true")
	if viper.GetEnvConfigBool("system.debug") {
		l.c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	l.Logger, err = l.c.Build()
	if err != nil {
		return err
	}
	l.Sugar = l.Logger.Sugar()
	return nil
}

func (l *Logger) newLogDir() {
	l.logFileDir = utils.GetPath() + "/runtime"
	Infof("LogDir ==> %s", l.logFileDir)
	if !utils.IsDirExists(l.logFileDir) {
		if err := utils.MkdirFile(l.logFileDir); err != nil {
			Errorf("utils.MkdirFile(%s),err:%+v.\n", l.logFileDir, err)
			l.logFileDir = ""
			l.err = err
		}
	}
}

func (l *Logger) newLogFile() {
	if l.err == nil && l.logFileDir != "" {
		var err error
		l.logFileName = l.logFileDir + "/" + time.Now().Format("2006-01-02") + ".log"
		if l.logFile != nil {
			l.logFile.Close()
		}
		l.logFile, err = os.OpenFile(l.logFileName, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			Error(err)
			l.logFile, err = os.Create(l.logFileName)
			if err != nil {
				l.logFileName = ""
				l.err = err
			}
		}
	}
}

// 检测是否跨天了,把记录记录到新的文件目录中
func (l *Logger) updateLogFile() {
	viper.C.SetDefault("system.saveDays", "7")
	saveDays := viper.GetEnvConfigFloat("system.saveDays")
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		select {
		case <-t.C:
			go deleteLog(l.logFileDir, saveDays) //删除修改时间在saveDays天前的文件

			l.newLogFile()

			if err := l.newZap(); err != nil {
				Errorf("l.newZap(),err:%+v.", err)
				continue
			}
		}
	}
}

func deleteLog(sourceDir string, saveDays float64) {
	filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Errorf("filepath.Walk.func,err:%+v.", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		t := time.Now().Sub(info.ModTime()).Hours()
		if t >= (saveDays-1)*24 {
			e := os.Remove(path)
			if e != nil {
				Errorf("filepath.Walk.os.Remove(%s),err:%+v.", path, err)
			}
		}
		return nil
	})
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, "2006-01-02 15:04:05.000", enc)
}

func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}
