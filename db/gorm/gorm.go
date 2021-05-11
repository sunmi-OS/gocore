package gorm

import (
	viper2 "github.com/sunmi-OS/gocore/conf/viper"
	retry2 "github.com/sunmi-OS/gocore/utils/retry"
	xlog2 "github.com/sunmi-OS/gocore/utils/xlog"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Client struct {
	gormMaps      sync.Map
	defaultDbName string
}

var _Gorm *Client

// 初始化Gorm
func NewDB(dbname string) (g *Client) {
	NewOrUpdateDB(dbname)
	return _Gorm
}

// SetDefaultName 设置默认DB Name
func SetDefaultName(dbName string) {
	_Gorm.defaultDbName = dbName
}

// Deprecated
// 推荐使用：NewOrUpdateDB
// Updata 更新Gorm集成新建
func UpdateDB(dbname string) error {
	return NewOrUpdateDB(dbname)
}

// NewOrUpdateDB 初始化或更新Gorm
func NewOrUpdateDB(dbname string) error {
	var (
		orm *gorm.DB
		err error
	)

	if _Gorm == nil {
		_Gorm = &Client{defaultDbName: defaultName}
	}

	// second: load gorm client
	oldGorm, _ := _Gorm.gormMaps.Load(dbname)

	// first: open new gorm client
	err = retry2.Retry(func() error {
		orm, err = openORM(dbname)
		if err != nil {
			xlog2.Errorf("UpdateDB(%s) error:%+v", dbname, err)
			return err
		}
		return nil
	}, 5, 3*time.Second)

	// 如果NEW异常直接panic如果是Update返回error
	if err != nil {
		if oldGorm == nil {
			panic(err)
		}
		return err
	}

	// third: delete old gorm client and store the new gorm client
	_Gorm.gormMaps.Delete(dbname)
	_Gorm.gormMaps.Store(dbname, orm)

	// fourth: if old client is not nil, delete and close connection
	if oldGorm != nil {
		oldGorm.(*gorm.DB).Close()
	}
	return nil
}

// GetORM 获取默认的Gorm实例
// 目前仅支持 不传 或者仅传一个 dbname
func GetORM(dbname ...string) *gorm.DB {
	name := _Gorm.defaultDbName
	if len(dbname) == 1 {
		name = dbname[0]
	}

	v, ok := _Gorm.gormMaps.Load(name)
	if ok {
		return v.(*gorm.DB)
	}
	return nil
}

func Close() {
	_Gorm.gormMaps.Range(func(dbName, orm interface{}) bool {
		xlog2.Warnf("close db %s", dbName)
		_Gorm.gormMaps.Delete(dbName)
		orm.(*gorm.DB).Close()
		return true
	})
}

//----------------------------以下是私有方法--------------------------------

// openORM 私有方法
func openORM(dbname string) (*gorm.DB, error) {
	//默认配置
	viper2.C.SetDefault(dbname, map[string]interface{}{
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
	dbHost := viper2.GetEnvConfig(dbname + ".dbHost")
	dbName := viper2.GetEnvConfig(dbname + ".dbName")
	dbUser := viper2.GetEnvConfig(dbname + ".dbUser")
	dbPasswd := viper2.GetEnvConfig(dbname + ".dbPasswd")
	dbPort := viper2.GetEnvConfig(dbname + ".dbPort")
	dbType := viper2.GetEnvConfig(dbname + ".dbType")
	dbDebug := viper2.GetEnvConfigBool(dbname + ".dbDebug")

	connectString := dbUser + ":" + dbPasswd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=true&loc=Local"
	orm, err := gorm.Open(dbType, connectString)
	if err != nil {
		return nil, err
	}

	// 连接池的空闲数大小
	orm.DB().SetMaxIdleConns(viper2.C.GetInt(dbname + ".dbIdleconns_max"))
	// 最大打开连接数
	orm.DB().SetMaxOpenConns(viper2.C.GetInt(dbname + ".dbOpenconns_max"))

	if dbDebug {
		// 开启Debug模式
		orm = orm.Debug()
	}
	return orm, nil
}
