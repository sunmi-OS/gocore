package main

import (
	"fmt"

	"github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/xlog"
)

type Machine struct {
	Mid int64  `gorm:"column:mid"`
	Msn string `gorm:"column:msn"`
}

func (m Machine) TableName() string {
	return "machine"
}

func main() {
	// 指定配置文件所在的目录和文件名称
	viper.NewConfig("config", "conf")

	gorm.NewDB("a")
	gorm.NewDB("b")
	gorm.NewDB("c")

	client := gorm.Gorm()
	err := client.NewOrUpdateDB("d")
	if err != nil {
		xlog.Errorf("NewOrUpdateDB(%s),error:%+v", "d", err)
	}

	var MC []Machine

	err = client.GetORM().Where("msn =  ?", "7102V04115500128").Find(&MC).Error

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(MC)

	client.Close()
}
