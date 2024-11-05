package orm

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/sunmi-OS/gocore/v2/utils/closes"
)

type Client struct {
	gormMaps      sync.Map
	defaultDbName string
}

var _Gorm *Client
var closeOnce sync.Once

// NewDB initialize db session
func NewDB(dbname string) (g *Client) {
	err := NewOrUpdateDB(dbname)
	if err != nil {
		return nil
	}
	closeOnce.Do(func() {
		closes.AddShutdown(closes.ModuleClose{
			Name:     "Gorm Close",
			Priority: closes.GormPriority,
			Func:     Close,
		})
	})
	return _Gorm
}

func SetDefaultName(dbName string) {
	_Gorm.defaultDbName = dbName
}

// NewOrUpdateDB initialize or update db session
func NewOrUpdateDB(dbname string) error {
	var (
		orm *gorm.DB
		err error
	)
	if _Gorm == nil {
		_Gorm = &Client{defaultDbName: defaultName}
	}
	// load gorm client
	oldGorm, _ := _Gorm.gormMaps.Load(dbname)
	// open new gorm client
	err = utils.Retry(func() error {
		orm, err = openORM(dbname)
		if err != nil {
			glog.ErrorF("UpdateDB(%s) error:%+v", dbname, err)
			return err
		}
		return nil
	}, 5, 3*time.Second)
	if err != nil {
		if oldGorm == nil {
			panic(err)
		}
		return err
	}
	_Gorm.gormMaps.Store(dbname, orm)
	// if old client is not nil, delete and close connection
	if oldGorm != nil {
		db, _ := oldGorm.(*gorm.DB).DB()
		if db != nil {
			err := db.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetORM 获取默认的Gorm实例
// 目前仅支持 不传 或者仅传一个 dbname
func GetORM(dbname ...string) *gorm.DB {
	name := _Gorm.defaultDbName
	if len(dbname) != 0 {
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
		glog.WarnF("close db %s", dbName)
		_Gorm.gormMaps.Delete(dbName)
		db, _ := orm.(*gorm.DB).DB()
		if db != nil {
			db.Close()
		}
		return true
	})
}

// openORM initialize db session
func openORM(dbname string) (*gorm.DB, error) {
	// default setting
	viper.C.SetDefault(dbname, map[string]interface{}{
		"Port":            3306,
		"MaxIdleConns":    20,
		"MaxOpenConns":    20,
		"Type":            "mysql",
		"Debug":           false,
		"MultiStatements": false,
	})
	dbHost := viper.GetEnvConfig(dbname + ".Host").String()
	if dbHost == "" {
		return nil, errors.New(fmt.Sprintf("the Host in the %s database configuration file is empty", dbname))
	}
	dbName := viper.GetEnvConfig(dbname + ".Name").String()
	if dbName == "" {
		return nil, errors.New(fmt.Sprintf("the Name in the %s database configuration file is empty", dbname))
	}
	dbUser := viper.GetEnvConfig(dbname + ".User").String()
	if dbUser == "" {
		return nil, errors.New(fmt.Sprintf("the User in the %s database configuration file is empty", dbname))
	}
	dbPasswd := viper.GetEnvConfig(dbname + ".Passwd").String()
	if dbPasswd == "" {
		return nil, errors.New(fmt.Sprintf("the Passwd in the %s database configuration file is empty", dbname))
	}
	dbPort := viper.GetEnvConfig(dbname + ".Port").String()
	dbType := viper.GetEnvConfig(dbname + ".Type").String()
	dbDebug := viper.GetEnvConfig(dbname + ".Debug").Bool()
	dbMulti := viper.GetEnvConfig(dbname + ".MultiStatements").Bool()

	dsn := dbUser + ":" + dbPasswd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=true&loc=Local"
	if dbMulti {
		dsn += "&multiStatements=true"
	}
	switch dbType {
	case "mysql":
		orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: glog.NewDBLogger(dbDebug), SkipDefaultTransaction: true})
		if err != nil {
			return nil, err
		}
		db, err := orm.DB()
		if err != nil {
			return nil, err
		}
		// 连接池的空闲数大小
		db.SetMaxIdleConns(viper.C.GetInt(dbname + ".MaxIdleConns"))
		// 最大打开连接数
		db.SetMaxOpenConns(viper.C.GetInt(dbname + ".MaxOpenConns"))
		return orm, nil
	default:
		return nil, fmt.Errorf("not support sql driver [%s]", dbType)
	}
}
