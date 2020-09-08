package gorm

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sunmi-OS/gocore/retry"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/xlog"
)

var (
	Gorm        sync.Map
	defaultName = "dbDefault"
)

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
	var (
		orm *gorm.DB
		err error
	)

	err = retry.Retry(func() error {
		orm, err = openORM(dbname)
		if err != nil {
			xlog.Errorf("NewDB(%s) error:%+v", dbname, err)
			return err
		}
		return nil
	}, 5, 3*time.Second)
	if err != nil || orm == nil {
		panic(err)
	}

	Gorm.Store(dbname, orm)
}

// 设置获取db的默认值
func SetDefaultName(dbname string) {
	defaultName = dbname
}

// 初始化Gorm
func UpdateDB(dbname string) error {
	var (
		orm *gorm.DB
		err error
	)

	// first: open new gorm client
	err = retry.Retry(func() error {
		orm, err = openORM(dbname)
		if err != nil {
			xlog.Errorf("UpdateDB(%s) error:%+v", dbname, err)
			return err
		}
		return nil
	}, 5, 3*time.Second)
	if err != nil {
		return err
	}

	// second: load gorm client
	v, _ := Gorm.Load(dbname)

	// third: delete old gorm client and store the new gorm client
	Gorm.Delete(dbname)
	Gorm.Store(dbname, orm)

	// fourth: if old client is not nil, delete and close connection
	if v != nil {
		v.(*gorm.DB).Close()
	}
	return nil
}

// Deprecated
// 通过名称获取Gorm实例
func GetORMByName(dbname string) *gorm.DB {
	v, ok := Gorm.Load(dbname)
	if ok {
		return v.(*gorm.DB)
	}
	return nil
}

// GetORM 获取默认的Gorm实例
// 目前仅支持 不传 或者仅传一个 dbname
func GetORM(dbname ...string) *gorm.DB {
	name := defaultName
	if len(dbname) == 1 {
		name = dbname[0]
	}
	v, ok := Gorm.Load(name)
	if ok {
		return v.(*gorm.DB)
	}
	return nil
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

	connectString := dbUser + ":" + dbPasswd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=true&loc=Local"
	orm, err := gorm.Open(dbType, connectString)
	if err != nil {
		return nil, err
	}

	// 连接池的空闲数大小
	orm.DB().SetMaxIdleConns(viper.C.GetInt(dbname + ".dbIdleconns_max"))
	// 最大打开连接数
	orm.DB().SetMaxOpenConns(viper.C.GetInt(dbname + ".dbOpenconns_max"))

	if dbDebug {
		// 开启Debug模式
		orm = orm.Debug()
	}
	return orm, nil
}
