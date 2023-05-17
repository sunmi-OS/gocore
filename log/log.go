package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sunmi-OS/gocore/utils"
	"github.com/sunmi-OS/gocore/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger  *zap.Logger
	Sugar   *zap.SugaredLogger
	logfile *os.File
	cfg     zap.Config
	once    sync.Once
)

type Option func(*zap.Config)

type OutPut int

const (
	ConsoleOutPut OutPut = iota + 1
	FileOutPut
	ConsoleAndFileOutPut
)

// InitLogger 初始化Log日志记录
func InitLogger(serviceName string, opt ...Option) {
	var err error
	cfg = zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	Logger, err = cfg.Build()
	if err != nil {
		panic("init logger error")
	}
	Sugar = Logger.Sugar()
	for _, option := range opt {
		option(&cfg)
	}
}

// SetOutputPath Set log output Path
func SetOutputPath(outputPath OutPut) Option {
	return func(c *zap.Config) {
		if outputPath == FileOutPut || outputPath == ConsoleAndFileOutPut {
			var err error
			if !utils.IsDirExists(utils.GetPath() + "/Runtime") {
				if mkdirerr := utils.MkdirFile(utils.GetPath() + "/Runtime"); mkdirerr != nil {
					panic(mkdirerr)
				}
			}
			filename := utils.GetPath() + "/Runtime/" + time.Now().Format("2006-01-02") + ".log"
			logfile, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				logfile, err = os.Create(filename)
				if err != nil {
					panic(err)
				}
			}
			if outputPath == FileOutPut {
				c.OutputPaths = []string{filename}
				c.ErrorOutputPaths = []string{filename}
			} else {
				c.OutputPaths = append(c.OutputPaths, filename)
				c.ErrorOutputPaths = append(c.ErrorOutputPaths, filename)
			}
			Logger, err = c.Build()
			if err != nil {
				panic("init logger error")
			}
			Sugar = Logger.Sugar()
			once.Do(func() {
				go updateLogFile()
			})
		}
	}
}

// SetLogLevel alters the logging level.
func SetLogLevel(l zapcore.Level) Option {
	return func(c *zap.Config) {
		var err error
		c.Level = zap.NewAtomicLevelAt(l)
		Logger, err = c.Build()
		if err != nil {
			panic("init logger error")
		}
		Sugar = Logger.Sugar()
	}
}

// 检测是否跨天了,把记录记录到新的文件目录中
func updateLogFile() {
	var err error
	viper.C.SetDefault("system.saveDays", "7")
	saveDays := viper.GetEnvConfigFloat("system.saveDays")
	logPath := utils.GetPath() + "/Runtime/"
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		select {
		case <-t.C:
			// 以下为定时执行的操作
			logfile.Close()
			go deleteLog(logPath, saveDays) // 删除修改时间在saveDays天前的文件
			filename := logPath + time.Now().Format("2006-01-02") + ".log"
			logfile, err = os.Create(logPath + time.Now().Format("2006-01-02") + ".log")
			if err != nil {
				fmt.Println(err)
			}
			var outPaths []string
			var errPaths []string
			for _, v := range cfg.OutputPaths {
				if v == "stdout" {
					// cfg.ErrorOutputPaths = []string{filename, "stderr"}
					outPaths = append(outPaths, v)
					errPaths = append(errPaths, "stderr")
				} else {
					outPaths = append(outPaths, filename)
					errPaths = append(errPaths, filename)
				}
			}
			cfg.OutputPaths = outPaths
			cfg.ErrorOutputPaths = errPaths
			Logger, err = cfg.Build()
			if err != nil {
				fmt.Println(err)
				continue
			}
			Sugar = Logger.Sugar()
		}
	}
}

func deleteLog(source string, saveDays float64) error {

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		t := time.Now().Sub(info.ModTime()).Hours()
		fmt.Println(path, t)
		if t >= (saveDays-1)*24 {
			e := os.Remove(path)
			if e != nil {
				fmt.Println(e)
			}
		}
		return err
	})
}
