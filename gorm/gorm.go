package gorm

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sunmi-OS/gocore/viper"
)

var Gorm sync.Map
var defaultName = "dbDefault"


var (
	// ErrRecordNotFound record not found error, happens when haven't find any matched data when looking up with a struct
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidSQL invalid SQL error, happens when you passed invalid SQL
	ErrInvalidSQL = gorm.ErrInvalidSQL
	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
	// ErrCantStartTransaction can't start transaction when you are trying to start one with `Begin`
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	// ErrUnaddressable unaddressable value
	ErrUnaddressable = gorm.ErrUnaddressable
)

// 初始化Gorm
func NewDB(dbname string) {

	var orm *gorm.DB
	var err error

	for orm, err = openORM(dbname); err != nil; {
		fmt.Println("Database connection exception! 5 seconds to retry")
		time.Sleep(5 * time.Second)
		orm, err = openORM(dbname)
	}

	Gorm.LoadOrStore(dbname, orm)
}

// 设置获取db的默认值
func SetDefaultName(dbname string) {
	defaultName = dbname
}

// 初始化Gorm
func UpdateDB(dbname string) error {

	v, _ := Gorm.Load(dbname)

	orm, err := openORM(dbname)

	Gorm.Delete(dbname)
	Gorm.LoadOrStore(dbname, orm)

	err = v.(*gorm.DB).Close()
	if err != nil {
		return err
	}

	return nil
}

// 通过名称获取Gorm实例
func GetORMByName(dbname string) *gorm.DB {

	v, _ := Gorm.Load(dbname)
	return v.(*gorm.DB)
}

// 获取默认的Gorm实例
func GetORM() *gorm.DB {

	v, _ := Gorm.Load(defaultName)
	return v.(*gorm.DB)
}

func openORM(dbname string) (*gorm.DB, error) {

	//默认配置
	viper.C.SetDefault(dbname, map[string]interface{}{
		"dbHost":          "127.0.0.1",
		"dbName":          "phalgo",
		"dbUser":          "root",
		"dbPasswd":        "",
		"dbPort":          3306,
		"dbIdleconns_max": 20,
		"dbOpenconns_max": 20,
		"dbType":          "mysql",
		"dbDebug":         false,
	})
	dbHost := viper.GetEnvConfig(dbname + ".dbHost")
	dbName := viper.GetEnvConfig(dbname + ".dbName")
	dbUser := viper.GetEnvConfig(dbname + ".dbUser")
	dbPasswd := viper.GetEnvConfig(dbname + ".dbPasswd")
	dbPort := viper.GetEnvConfig(dbname + ".dbPort")
	dbType := viper.GetEnvConfig(dbname + ".dbType")
	dbDebug := viper.GetEnvConfigBool(dbname + ".dbDebug")

	connectString := dbUser + ":" + dbPasswd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&parseTime=true&loc=Local"

	orm, err := gorm.Open(dbType, connectString)

	if err != nil {
		return nil, err
	}

	//连接池的空闲数大小
	orm.DB().SetMaxIdleConns(viper.C.GetInt(dbname + ".dbIdleconns_max"))
	//最大打开连接数
	orm.DB().SetMaxOpenConns(viper.C.GetInt(dbname + ".dbOpenconns_max"))

	if dbDebug {
		// 开启Debug模式
		orm = orm.Debug()
	}

	return orm, err
}
