package main

import (
	"fmt"
	"github.com/sunmi-OS/gocore/example/nacos/config"
	"github.com/sunmi-OS/gocore/gorm"
	"github.com/sunmi-OS/gocore/viper"
	"time"

	"github.com/sunmi-OS/gocore/nacos"
)

type App struct {
	Description string `gorm:"description"`
}

func main() {

	config.InitNacos("dev")

	nacos.SetDataIds("DEFAULT_GROUP", "adb")
	nacos.SetDataIds("pay", "test")


	nacos.SetCallBackFunc("DEFAULT_GROUP", "adb", func(namespace, group, dataId, data string) {

		err := gorm.UpdateDB("remotemanageDB")
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	nacos.NacosToViper()



	s := viper.C.GetString("remotemanageDB.dbHost")

	fmt.Println(s)

	s = viper.C.GetString("redisDB.remote_control")

	fmt.Println(s)

	gorm.NewDB("remotemanageDB")

	i := 0
	for {

		orm1 := gorm.GetORMByName("remotemanageDB")

		app := App{}

		err := orm1.Raw("select description from app").Find(&app).Error

		fmt.Println(app)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Print("ping ok", i)
			i++
		}

		time.Sleep(time.Second * 1)

	}

	time.Sleep(time.Second * 1000)

}
