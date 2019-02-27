package viper

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/sunmi-OS/gocore/utils"
)

var C *viper.Viper

// 初始化配置文件
// filePath 配置文件路径
// fileName 配置文件名称(不需要文件后缀)
func NewConfig(filePath string, fileName string) {

	C = viper.New()

	C.SetConfigName(fileName)
	//filePath支持相对路径和绝对路径 etc:"/a/b" "b" "./b"
	if (filePath[:1] != "/") {
		C.AddConfigPath(path.Join(utils.GetPath(), filePath))
	} else {
		C.AddConfigPath(filePath)
	}

	C.WatchConfig()
	
	// 找到并读取配置文件并且 处理错误读取配置文件
	if err := C.ReadInConfig(); err != nil {
		panic(err)
	}

}


// 获取配置文件优先获取环境变量(返回string类型)
func GetEnvConfig(key string) string {

	// 转大写 . 转 _ 获取环境变量判断是否存在(存在直接返回,不存在使用viper配置)
	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return env
	}

	return C.GetString(key)
}

// 获取配置文件优先获取环境变量(返回int类型)
func GetEnvConfigInt(key string) int64 {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToInt64(env)
	}

	return C.GetInt64(key)
}

// 获取配置文件优先获取环境变量(返回Float类型)
func GetEnvConfigFloat(key string) float64 {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToFloat64(env)
	}

	return C.GetFloat64(key)
}

// 获取配置文件优先获取环境变量(返回Bool类型)
func GetEnvConfigBool(key string) bool {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToBool(env)
	}

	return C.GetBool(key)
}


