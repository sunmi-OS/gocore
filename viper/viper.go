package viper

import (
	"path"
	"strings"
	"os"

	"github.com/spf13/viper"
	"github.com/sunmi-OS/gocore/utils"
	"github.com/spf13/cast"
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

// 转大写 . 转 _ 获取环境变量判断是否存在(存在直接返回,不存在使用viper配置)
// 使用获取配置需要Docker来获取环境变量的场景
func GetEnvConfig(key string) string {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return env
	}

	return C.GetString(key)
}

func GetEnvConfigInt(key string) int64 {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToInt64(env)
	}

	return C.GetInt64(key)
}

func GetEnvConfigFloat(key string) float64 {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToFloat64(env)
	}

	return C.GetFloat64(key)
}

func GetEnvConfigBool(key string) bool {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToBool(env)
	}

	return C.GetBool(key)
}


