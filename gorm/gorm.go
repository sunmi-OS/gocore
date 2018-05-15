package gorm

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sunmi-OS/gocore/viper"
)

var Gorm map[string]*gorm.DB

func init() {
	Gorm = make(map[string]*gorm.DB)
}

// 初始化Gorm
func NewDB(dbname string) {

	var orm *gorm.DB
	var err error

	//默认配置
	viper.C.SetDefault(dbname, map[string]interface{}{
		"dbHost":          "127.0.0.1",
		"dbName":          "phalgo",
		"dbUser":          "root",
		"dbPasswd":        "",
		"dbPort":          3306,
		"dbIdleconns_max": 0,
		"dbOpenconns_max": 20,
		"dbType":          "mysql",
	})
	dbHost := viper.C.GetString(dbname + ".dbHost")
	dbName := viper.C.GetString(dbname + ".dbName")
	dbUser := viper.C.GetString(dbname + ".dbUser")
	dbPasswd := viper.C.GetString(dbname + ".dbPasswd")
	dbPort := viper.C.GetString(dbname + ".dbPort")
	dbType := viper.C.GetString(dbname + ".dbType")

	connectString := dbUser + ":" + dbPasswd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"
	//开启sql调试模式
	//GDB.LogMode(true)

	for orm, err = gorm.Open(dbType, connectString); err != nil; {
		fmt.Println("数据库连接异常! 5秒重试")
		time.Sleep(5 * time.Second)
		orm, err = gorm.Open(dbType, connectString)
	}

	//连接池的空闲数大小
	orm.DB().SetMaxIdleConns(viper.C.GetInt(dbname + ".idleconns_max"))
	//最大打开连接数
	orm.DB().SetMaxIdleConns(viper.C.GetInt(dbname + ".openconns_max"))
	Gorm[dbname] = orm
	//defer Gorm[dbname].Close()
}

// 通过名称获取Gorm实例
func GetORMByName(dbname string) *gorm.DB {

	return Gorm[dbname]
}

// 获取默认的Gorm实例
func GetORM() *gorm.DB {

	return Gorm["dbDefault"]
}
