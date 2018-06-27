package main

import (
	"fmt"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/log"
	"github.com/sunmi-OS/gocore/hbase"
	"os"
	"go.uber.org/zap"
)

func main() {

	// 初始化配置文件
	viper.NewConfig("config", "conf")
	// 初始化日志库
	log.InitLogger("example-Hbase")

	err := hbase.NewHbase()

	if err != nil {
		fmt.Println("连接失败,错误原因:", err.Error())
		os.Exit(0)
	}

	talbeName := []byte("sunmi6")
	rowkey := []byte("10009")

	cvarr := []*hbase.TColumnValue{
		{Family: []byte("info"), Qualifier: []byte("age"), Value: []byte("23")},
		{Family: []byte("info"), Qualifier: []byte("name"), Value: []byte("wenzhenxi")},
		{Family: []byte("info"), Qualifier: []byte("sex"), Value: []byte("n")},
	}

	// 写入数据
	htput := hbase.TPut{Row: rowkey, ColumnValues: cvarr}
	err = hbase.HbaseClinet.Put(talbeName, &htput)
	if err != nil {
		log.Sugar.Infow("写入Hbase数据失败,原因:" + err.Error())
	} else {
		log.Sugar.Infow("写入Hbase数据成功!")
	}
	// 获取数据
	objs, err := hbase.HbaseClinet.Get(talbeName, &hbase.TGet{Row: rowkey})
	if err != nil {
		log.Sugar.Infow("获取数据失败,原因:" + err.Error())
	} else {
		log.Sugar.Infow("获取Hbase数据成功!", zap.Any("objs", objs))
	}
	// 删除数据
	err = hbase.HbaseClinet.Delete(talbeName, &hbase.TDelete{Row: rowkey})
	if err != nil {
		log.Sugar.Infow("删除Hbase数据失败,原因:" + err.Error())
	} else {
		log.Sugar.Infow("删除Hbase数据成功!")
	}
	// 判读数据是否存在
	b, err := hbase.HbaseClinet.Exists(talbeName, &hbase.TGet{Row: rowkey})
	if err != nil {
		log.Sugar.Infow("获取数据失败,原因:" + err.Error())
	} else {
		if b {
			log.Sugar.Infow("数据存在")
		} else {
			log.Sugar.Infow("数据不存在")
		}
	}
	
}
