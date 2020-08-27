package xlog

import (
	"log"
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
	Logger      *zap.Logger
	Sugar       *zap.SugaredLogger
	cfg         zap.Config
	logFile     *os.File
	logFileDir  string
	logFileName string
	once        sync.Once
}

var (
	logger = &Logger{}
)

func init() {
	logger.once.Do(func() {
		logger.init()
	})
}

func (l *Logger) init() {
	// new log file dir
	l.logFileDir = utils.GetPath() + "/runtime"
	if !utils.IsDirExists(l.logFileDir) {
		if err := utils.MkdirFile(l.logFileDir); err != nil {
			log.Printf("utils.MkdirFile(%s),err:%+v.\n", l.logFileDir, err)
			// todo 确认是否return
			return
		}
	}
	// new log file
	l.newLogFile()
	// init zap
	l.initZap()
	// updateLogFile Loop
	go l.updateLogFile()
}

func (l *Logger) newLogFile() (err error) {
	l.logFileName = l.logFileDir + time.Now().Format("2006-01-02") + ".log"
	if l.logFile != nil {
		l.logFile.Close()
	}
	l.logFile, err = os.OpenFile(l.logFileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		l.logFile, err = os.Create(l.logFileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) initZap() (err error) {
	l.cfg = zap.NewProductionConfig()
	l.cfg.OutputPaths = []string{l.logFileName, "stdout"}
	l.cfg.ErrorOutputPaths = []string{l.logFileName, "stderr"}
	l.cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	viper.C.SetDefault("system.debug", "true")
	if viper.GetEnvConfigBool("system.debug") {
		l.cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	l.Logger, err = l.cfg.Build()
	if err != nil {
		return err
	}
	l.Sugar = l.Logger.Sugar()
	return nil
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

			if err := l.newLogFile(); err != nil {
				log.Printf("l.newLogFile(%s),err:%+v.\n", l.logFileName, err)
				continue
			}

			if err := l.initZap(); err != nil {
				log.Printf("l.initZap(),err:%+v.\n", err)
				continue
			}
		}
	}
}

func deleteLog(sourceDir string, saveDays float64) {
	filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("filepath.Walk.func,err:%+v.\n", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		t := time.Now().Sub(info.ModTime()).Hours()
		if t >= (saveDays-1)*24 {
			e := os.Remove(path)
			if e != nil {
				log.Printf("filepath.Walk.os.Remove(%s),err:%+v.\n", path, err)
			}
		}
		return nil
	})
}
