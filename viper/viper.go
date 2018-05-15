package viper

import (
	"path"

	"github.com/spf13/viper"
	"github.com/sunmi-OS/gocore/utils"
)

var C *viper.Viper

//初始化配置文件
func NewConfig(filePath string, fileName string) {

	C = viper.New()
	C.WatchConfig()
	C.SetConfigName(fileName)
	//filePath支持相对路径和绝对路径 etc:"/a/b" "b" "./b"
	if (filePath[:1] != "/") {
		C.AddConfigPath(path.Join(utils.GetPath(), filePath))
	} else {
		C.AddConfigPath(filePath)
	}

	// 找到并读取配置文件并且 处理错误读取配置文件
	if err := C.ReadInConfig(); err != nil {
		panic(err)
	}

}


