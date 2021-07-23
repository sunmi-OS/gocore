package file

import (
	"os"
	"path/filepath"
	"strings"
)

// ApiRoot 当前项目根目录
var ApiRoot string

// GetPath 获取项目路径
func GetPath() string {
	if ApiRoot != "" {
		return ApiRoot
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		print(err.Error())
	}
	ApiRoot = strings.Replace(dir, "\\", "/", -1)
	return ApiRoot
}

// CheckDir 判断文件目录否存在
func CheckDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

// MkdirDir 创建文件夹,支持x/a/a  多层级
func MkdirDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// RemoveDir 删除文件
func RemoveDir(filePath string) error {
	return os.RemoveAll(filePath)
}
