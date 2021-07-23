package file

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"reflect"
	"strings"
)

// CheckFile 判断文件是否存在  存在返回 true 不存在返回false
func CheckFile(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetFileSize 获取文件大小
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	size := fileInfo.Size()
	return size, nil
}

// UploadFile 创建文件并生产目录
func UploadFile(file *multipart.FileHeader, path string) (string, error) {
	if reflect.ValueOf(file).IsNil() || !reflect.ValueOf(file).IsValid() {
		return "", errors.New("invalid memory address or nil pointer dereference")
	}
	src, err := file.Open()
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Println(err)
		}
	}(src)
	if err != nil {
		return "", err
	}

	err = MkdirDir(path)
	if err != nil {
		return "", err
	}
	filename := strings.Replace(file.Filename, " ", "", -1)
	filename = strings.Replace(filename, "\n", "", -1)
	dst, err := os.Create(path + filename)
	if err != nil {
		return "", err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Println(err)
		}
	}(dst)

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return filename, nil
}
