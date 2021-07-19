package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	viper2 "github.com/sunmi-OS/gocore/v2/conf/viper"

	"github.com/sunmi-OS/gocore/v2/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger  *zap.Logger
	Sugar   *zap.SugaredLogger
	logfile *os.File
	cfg     zap.Config
)

// 初始化Log日志记录
func InitLogger(serviceName string) {
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
	cfg.OutputPaths = []string{filename, "stdout"}
	cfg.ErrorOutputPaths = []string{filename, "stderr"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	viper2.C.SetDefault("system.debug", "true")
	if viper2.GetEnvConfigBool("system.debug") {
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
	viper2.C.SetDefault("system.saveDays", "7")
	saveDays := viper2.GetEnvConfigFloat("system.saveDays")
	logPath := utils.GetPath() + "/Runtime/"
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
			go deleteLog(logPath, saveDays) //删除修改时间在saveDays天前的文件
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
