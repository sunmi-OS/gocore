package log

import (
	"fmt"
	"os"
	"time"

	"github.com/sunmi-OS/gocore/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/sunmi-OS/gocore/viper"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger
var logfile *os.File
var cfg zap.Config

// 初始化Log日志记录
func InitLogger(serviceaName string) {
	var err error

	if !utils.IsDirExists(utils.GetPath() + "/Runtime") {
		if mkdirerr := utils.MkdirFile(utils.GetPath() + "/Runtime"); mkdirerr != nil {
			fmt.Println(mkdirerr)
		}
	}

	filename := utils.GetPath() + "/Runtime/" + time.Now().Format("2006-01-02") + ".log"
	logfile, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logfile, err = os.Create(filename)
		if err != nil {
			fmt.Println(err)
		}
	}

	cfg = zap.NewProductionConfig()
	cfg.OutputPaths = []string{filename, "stderr"}
	cfg.ErrorOutputPaths = []string{filename, "stderr"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	viper.C.SetDefault("system.debug", "true")
	if viper.GetEnvConfigBool("system.debug") {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	Logger, err = cfg.Build()

	if err != nil {
		fmt.Println(err)
	}
	Sugar = Logger.Sugar()

	go updateLogFile()
}

// 检测是否跨天了,把记录记录到新的文件目录中
func updateLogFile() {
	var err error
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		select {
		case <-t.C:
			//以下为定时执行的操作
			logfile.Close()
			filename := utils.GetPath() + "/Runtime/" + time.Now().Format("2006-01-02") + ".log"
			logfile, err = os.Create(utils.GetPath() + "/Runtime/" + time.Now().Format("2006-01-02") + ".log")
			if err != nil {
				fmt.Println(err)
			}
			cfg.ErrorOutputPaths = []string{filename}
			cfg.OutputPaths = []string{filename}
			Logger, err = cfg.Build()
			if err != nil {
				fmt.Println(err)
				continue
			}
			Sugar = Logger.Sugar()
		}
	}
}
