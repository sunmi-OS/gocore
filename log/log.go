package log

import (
	"fmt"
	"os"
	"time"

	"github.com/sunmi-OS/gocore/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var day string
var logfile *os.File
var isDebug bool
var cfg zap.Config

// 初始化Log日志记录
func InitLogger(serviceaName string, debug bool) {
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
	cfg = zap.Config{
		Encoding:         "json",
		OutputPaths:      []string{filename},
		ErrorOutputPaths: []string{filename},
		Level:            zap.NewAtomicLevel(),
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		InitialFields:    map[string]interface{}{"service": serviceaName},
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	Logger, err = cfg.Build()
	if err != nil {
		fmt.Println(err)
	}
	Logger.Info("logger初始化成功")
	isDebug = debug
	day = time.Now().Format("02")

}

// 检测是否跨天了,把记录记录到新的文件目录中
func updateLogFile() {
	var err error
	day2 := time.Now().Format("02")
	if day2 != day {
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
			return
		}
		day = day2
	}
}

// 记录Debug信息
func LogDebug(msg string, fields ...zap.Field) {
	if isDebug == false {
		return
	}
	updateLogFile()
	Logger.Debug(msg, fields...)
}

// 记录Info信息
func LogInfo(msg string, fields ...zap.Field) {
	updateLogFile()
	Logger.Info(msg, fields...)
}

// 记录Error信息
func LogError(msg string, fields ...zapcore.Field) {
	updateLogFile()
	Logger.Error(msg, fields...)
}
