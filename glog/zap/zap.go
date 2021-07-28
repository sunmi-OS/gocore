package zap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/martian/log"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/utils/file"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger  *zap.Logger
	Sugar   *zap.SugaredLogger
	logfile *os.File
	cfg     zap.Config
)

func init() {
	var err error
	cfg = zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	Logger, err = cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Errorf("l.initZap(),err:%+v", err)
		return
	}
	Sugar = Logger.Sugar()
}

func SetLocLevel(logLevel string) {

	switch logLevel {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	Logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Errorf("l.initZap(),err:%+v", err)
		return
	}
	Sugar = Logger.Sugar()
}

func InitFileLog() {
	var err error

	fmt.Println(file.GetPath())
	if !file.CheckDir(file.GetPath() + "/Runtime") {
		if err := file.MkdirDir(file.GetPath() + "/Runtime"); err != nil {
			log.Errorf("l.initZap(),err:%+v", err)
		}
	}
	filename := file.GetPath() + "/Runtime/" + time.Now().Format("2006-01-02") + ".log"
	logfile, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logfile, err = os.Create(filename)
		if err != nil {
			fmt.Println(err)
		}
	}
	cfg.OutputPaths = []string{filename, "stdout"}
	cfg.ErrorOutputPaths = []string{filename, "stderr"}
	SetLocLevel(viper.GetEnvConfig("log.level").String())
	go updateLogFile()
}

// updateLogFile 检测是否跨天了,把记录记录到新的文件目录中
func updateLogFile() {
	var err error
	viper.C.SetDefault("log.saveDays", "7")
	saveDays := viper.GetEnvConfig("system.saveDays").Float64()
	logPath := file.GetPath() + "/Runtime/"
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时执行的操作
		logfile.Close()
		go deleteLog(logPath, saveDays)
		filename := logPath + time.Now().Format("2006-01-02") + ".log"
		logfile, err = os.Create(logPath + time.Now().Format("2006-01-02") + ".log")
		if err != nil {
			fmt.Println(err)
		}
		cfg.ErrorOutputPaths = []string{filename, "stderr"}
		cfg.OutputPaths = []string{filename, "stdout"}
		Logger, err = cfg.Build()
		if err != nil {
			fmt.Println(err)
			continue
		}
		Sugar = Logger.Sugar()
	}
}

// deleteLog 删除修改时间在saveDays天前的文件
func deleteLog(source string, saveDays float64) {
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && !strings.HasSuffix(info.Name(), ".log") {
			return nil
		}
		t := time.Since(info.ModTime()).Hours()
		if t >= (saveDays-1)*24 {
			e := os.Remove(path)
			if e != nil {
				fmt.Println(e)
			}
		}
		return err
	})
	if err != nil {
		return
	}
}
