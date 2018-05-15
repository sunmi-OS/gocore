package log

import (
	"os"
	"time"
	"fmt"
	
	"github.com/Sirupsen/logrus"
	"github.com/sunmi-OS/gocore/utils"
)

var LogS *logrus.Logger
var day string
var logfile *os.File

// 初始化Log日志记录
func init() {
	var err error
	LogS = logrus.New()
	LogS.Formatter = new(logrus.JSONFormatter)
	LogS.Level = logrus.DebugLevel

	if !utils.IsDirExists(utils.GetPath() + "/Runtime") {
		if mkdirerr := utils.MkdirFile(utils.GetPath() + "/Runtime"); mkdirerr != nil {
			fmt.Println(mkdirerr)
		}
	}

	logfile, err = os.OpenFile(utils.GetPath()+"/Runtime/"+time.Now().Format("2006-01-02")+".log", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logfile, err = os.Create(utils.GetPath() + "/Runtime/" + time.Now().Format("2006-01-02") + ".log")
		if err != nil {
			fmt.Println(err)
		}
	}
	LogS.Out = logfile
	day = time.Now().Format("02")

}

// 检测是否跨天了,把记录记录到新的文件目录中
func updateLogFile() {
	var err error
	day2 := time.Now().Format("02")
	if day2 != day {
		logfile.Close()
		logfile, err = os.Create(utils.GetPath() + "/Runtime/" + time.Now().Format("2006-01-02") + ".log")
		if err != nil {
			fmt.Println(err)
		}
		LogS.Out = logfile
		day = day2
	}
}

// 记录Debug信息
func LogDebug(str ...interface{}) {
	updateLogFile()
	LogS.Debug(str)
}

// 记录Info信息
func LogInfo(str ...interface{}) {
	updateLogFile()
	LogS.Info(str)
}

// 记录Error信息
func LogError(str ...interface{}) {
	updateLogFile()
	LogS.Error(str)
}
