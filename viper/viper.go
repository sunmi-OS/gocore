package viper

import (
	"bytes"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/sunmi-OS/gocore/utils"
)

type Viper struct {
	C *viper.Viper
}

var multipleViper sync.Map
var C = viper.New()

func NewConfigToToml(configs string) {

	C.SetConfigType("toml")
	err := C.ReadConfig(bytes.NewBuffer([]byte(configs)))
	if err != nil {
		print(err)
	}
}

func MerageConfigToToml(configs string) {
	C.SetConfigType("toml")
	err := C.MergeConfig(bytes.NewBuffer([]byte(configs)))
	if err != nil {
		print(err)
	}
}

// 初始化配置文件
// filePath 配置文件路径
// fileName 配置文件名称(不需要文件后缀)
func NewConfig(filePath string, fileName string) {
	C = newConfig(filePath, fileName).C
}

func newConfig(filePath string, fileName string) *Viper {

	v := viper.New()

	v.SetConfigName(fileName)
	//filePath支持相对路径和绝对路径 etc:"/a/b" "b" "./b"
	if filePath[:1] != "/" {
		v.AddConfigPath(path.Join(utils.GetPath(), filePath))
	} else {
		v.AddConfigPath(filePath)
	}

	v.WatchConfig()

	// 找到并读取配置文件并且 处理错误读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return &Viper{v}
}

func BuildVipers(filePath string, fileName ...string) {
	for _, v := range fileName {
		_, found := multipleViper.Load(v)
		if !found { //can not remap
			A := newConfig(filePath, v)
			multipleViper.Store(v, A)
		}
	}
}

func LoadViperByFilename(filename string) *Viper {
	value, _ := multipleViper.Load(filename)
	if value == nil {
		return nil
	} else {
		return value.(*Viper)
	}
}

// 获取配置文件优先获取环境变量(返回string类型)
func (V *Viper) GetEnvConfig(key string) string {

	// 转大写 . 转 _ 获取环境变量判断是否存在(存在直接返回,不存在使用viper配置)
	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return env
	}

	return V.C.GetString(key)
}

// 获取配置文件优先获取环境变量(返回int类型)
func (V *Viper) GetEnvConfigInt(key string) int64 {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToInt64(env)
	}

	return V.C.GetInt64(key)
}

// 获取配置文件优先获取环境变量(返回Float类型)
func (V *Viper) GetEnvConfigFloat(key string) float64 {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToFloat64(env)
	}

	return V.C.GetFloat64(key)
}

// 获取配置文件优先获取环境变量(返回Bool类型)
func (V *Viper) GetEnvConfigBool(key string) bool {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return cast.ToBool(env)
	}

	return V.C.GetBool(key)
}

func (V *Viper) GetEnvConfigStringSlice(key string) []string {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return strings.Split(env, ",")
	}

	return V.C.GetStringSlice(key)
}

func (V *Viper) GetEnvConfigCastInt(key string) int {
	return int(V.GetEnvConfigInt(key))
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

func GetEnvConfigStringSlice(key string) []string {

	env := os.Getenv(strings.Replace(strings.ToUpper(key), ".", "_", -1))
	if env != "" {
		return strings.Split(env, ",")
	}

	return C.GetStringSlice(key)
}

func GetEnvConfigCastInt(key string) int {
	return int(GetEnvConfigInt(key))
}
