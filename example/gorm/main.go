package main

import (
	"fmt"
	"github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/viper"
)

type Machine struct {
	MId int64  `gorm:"column:mId"`
	Msn string `gorm:"column:msn"`
}

func (m Machine) TableName() string {
	return "machine"
}

func main() {
	// 指定配置文件所在的目录和文件名称
	viper.NewConfig("config", "conf")

	gorm.NewDB("dbDefault")

	MC := []Machine{}

	err := gorm.GetORM().Where("msn =  ?", "7102V04115500128").Find(&MC).Error

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(MC)

}
