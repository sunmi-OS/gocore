package utils

import (
	"path/filepath"
	"os"
	"strings"
)

// 当前项目根目录
var API_ROOT string

// 获取项目路径
func GetPath() string {

	if API_ROOT != "" {
		return API_ROOT
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		print(err.Error())
	}

	API_ROOT = strings.Replace(dir, "\\", "/", -1)
	return API_ROOT
}

// 判断文件目录否存在
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

}

// 创建文件
func MkdirFile(path string) error {

	err := os.Mkdir(path, os.ModePerm) //在当前目录下生成md目录
	if err != nil {
		return err
	}
	return nil
}
