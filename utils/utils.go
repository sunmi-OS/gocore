package utils

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"strings"
	"time"
)

// 返回当前时间
func GetDate() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05")
}

// 获取当前系统环境
func GetRunTime() string {
	//获取系统环境变量
	RUN_TIME := os.Getenv("RUN_TIME")
	if RUN_TIME == "" {
		fmt.Println("No RUN_TIME Can't start")
	}
	return RUN_TIME
}

// MD5 加密字符串
func GetMD5(plainText string) string {
	h := md5.New()
	h.Write([]byte(plainText))
	return hex.EncodeToString(h.Sum(nil))
}

//计算文件的md5，适用于本地文件计算
func GetMd5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(md5hash.Sum(nil)), nil
}

//从流中直接读取数据计算md5 并返回流的副本，不能用于计算大文件流否则内存占用很大
//@return io.Reader @params file的副本
func GetMd52(file io.Reader) (io.Reader, string, error) {
	var b bytes.Buffer
	md5hash := md5.New()
	if _, err := io.Copy(&b, io.TeeReader(file, md5hash)); err != nil {
		return nil, "", err
	}
	return &b, hex.EncodeToString(md5hash.Sum(nil)), nil
}

//解压
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func UploadFile(file *multipart.FileHeader, path string) (string, error) {
	if reflect.ValueOf(file).IsNil() || !reflect.ValueOf(file).IsValid() {
		return "", errors.New("invalid memory address or nil pointer dereference")
	}
	src, err := file.Open()
	defer src.Close()
	if err != nil {
		return "", err
	}
	err = MkDir(path)
	if err != nil {
		return "", err
	}
	// Destination
	// 去除空格
	filename := strings.Replace(file.Filename, " ", "", -1)
	// 去除换行符
	filename = strings.Replace(filename, "\n", "", -1)

	dst, err := os.Create(path + filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return filename, nil
}

func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	//文件大小
	fsize := fileInfo.Size()
	return fsize, nil
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
