package orm

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sunmi-OS/gocore/retry"
	"github.com/sunmi-OS/gocore/viper"
	"github.com/sunmi-OS/gocore/xlog"
)

type Client struct {
	maps          sync.Map
	defaultDbName string
}

var _Gorm *Client

func Gorm() *Client {
	return _Gorm
}

// 初始化Gorm
func NewGorm(dbname string) {
	var (
		orm *gorm.DB
		err error
	)
	if _Gorm == nil {
		_Gorm = &Client{defaultDbName: defaultName}
	}

	// openORM
	err = retry.Retry(func() error {
		orm, err = openORM(dbname)
		if err != nil {
			xlog.Errorf("NewGorm(%s) error:%+v", dbname, err)
			return err
		}
		return nil
	}, 5, 3*time.Second)
	if err != nil || orm == nil {
		panic(err)
	}

	// store db client
	_Gorm.maps.Store(dbname, orm)
}

// SetDefaultName 设置默认DB Name
func (c *Client) SetDefaultName(dbName string) {
	c.defaultDbName = dbName
}

// NewOrUpdateDB 初始化或更新Gorm
func (c *Client) NewOrUpdateDB(dbname string) error {
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
	v, _ := c.maps.Load(dbname)

	// third: delete old gorm client and store the new gorm client
	c.maps.Delete(dbname)
	c.maps.Store(dbname, orm)

	// fourth: if old client is not nil, delete and close connection
	if v != nil {
		v.(*gorm.DB).Close()
	}
	return nil
}

// GetORM 获取默认的Gorm实例
// 目前仅支持 不传 或者仅传一个 dbname
func (c *Client) GetORM(dbname ...string) *gorm.DB {
	name := c.defaultDbName
	if len(dbname) == 1 {
		name = dbname[0]
	}

	v, ok := c.maps.Load(name)
	if ok {
		return v.(*gorm.DB)
	}
	return nil
}

func (c *Client) Close() {
	c.maps.Range(func(dbName, orm interface{}) bool {
		xlog.Warnf("close db %s", dbName)
		c.maps.Delete(dbName)
		orm.(*gorm.DB).Close()
		return true
	})
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
