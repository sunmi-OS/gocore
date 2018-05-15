//	PhalGo-Config
//	使用spf13大神的viper配置文件获取工具作为phalgo的配置文件工具
//	喵了个咪 <wenzhenxi@vip.qq.com> 2016/5/11
//  依赖情况:
//          "github.com/spf13/viper"

package viper

import (
	"github.com/spf13/viper"
	"path"
	"BITU-service/core/base"
)

var C *viper.Viper

//初始化配置文件
func NewConfig(filePath string, fileName string) {

	C = viper.New()
	C.WatchConfig()
	C.SetConfigName(fileName)
	//filePath支持相对路径和绝对路径 etc:"/a/b" "b" "./b"
	if (filePath[:1] != "/") {
		C.AddConfigPath(path.Join(base.GetPath(), filePath))
	} else {
		C.AddConfigPath(filePath)
	}

	// 找到并读取配置文件并且 处理错误读取配置文件
	if err := C.ReadInConfig(); err != nil {
		panic(err)
	}

}


