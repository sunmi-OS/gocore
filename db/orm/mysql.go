package orm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/sunmi-OS/gocore/v2/utils/closes"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SetDbConn set mysql connect, if old connect exist, then overwrite it and close old connect
// if use this function, make sure the connection right is the responsibility of the user
func SetDbConn(name string, orm *gorm.DB) {
	if orm == nil {
		return
	}

	if _Gorm == nil {
		_Gorm = &Client{defaultDbName: defaultName}
		_Gorm.gormMaps.Store(name, orm)
		return
	}

	// if have old connection, close it
	oldConn, _ := _Gorm.gormMaps.Load(name)
	_Gorm.gormMaps.Store(name, orm)

	if oldConn != nil {
		db, _ := oldConn.(*gorm.DB).DB()
		if db != nil {
			db.Close()
		}
	}
}

// InitDbConn init db connection
func InitDbConn(dbname string, opts ...Options) *Client {
	// if opts does not exist, then use NewDB
	if opts == nil || len(opts) == 0 {
		return NewDB(dbname)
	}

	info := NewConnInfo(dbname, opts...)
	err := info.checkConnParams()
	if err != nil {
		glog.ErrorF("check %s connection params failed, err:%s", info.Type, err.Error())
		panic(err)
	}

	err = info.NewOrUpdateDB(dbname)
	if err != nil {
		glog.ErrorF("init db conn failed, err:%s", err.Error())
	}

	closeOnce.Do(func() {
		closes.AddShutdown(closes.ModuleClose{
			Name:     "Gorm Close " + dbname,
			Priority: closes.GormPriority,
			Func:     Close,
		})
	})

	return _Gorm
}

// checkConnParams if Dsn not null, use it, or Use the Host, Port, Username, Password generate it
func (i *ConnInfo) checkConnParams() error {
	if i.Dsn != "" && i.Database == "" {
		// parse dsn for database
		dsnSplits := strings.Split(i.Dsn, "?")
		if len(dsnSplits) == 1 {
			i.Dsn += "?charset=utf8mb4&parseTime=True&loc=Local"
		}

		connInfos := strings.Split(dsnSplits[0], "/")
		if len(connInfos) != 2 {
			return errors.New("dsn format error,dns: " + i.Dsn)
		}

		i.Database = connInfos[1]
		return nil
	}

	if i.Username == "" {
		return fmt.Errorf("username is empty")
	}

	if i.Password == "" {
		return fmt.Errorf("password is empty")
	}
	if i.Host == "" {
		return fmt.Errorf("host is empty")
	}
	if i.Port == 0 {
		return fmt.Errorf("port is empty")
	}
	if i.Database == "" {
		return fmt.Errorf("database is empty")
	}

	i.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		i.Username, i.Password, i.Host, i.Port, i.Database)
	if i.MultiStatements {
		i.Dsn += "&multiStatements=true"
	}

	return nil
}

// NewOrUpdateDB new db connect, if success, update _Gorm
func (i *ConnInfo) NewOrUpdateDB(dbname string) error {
	if _Gorm == nil {
		_Gorm = &Client{defaultDbName: defaultName}
	}
	if dbname == "" {
		dbname = i.Database
	}
	oldGorm, _ := _Gorm.gormMaps.Load(dbname)
	err := utils.Retry(func() error {
		orm, err := i.connectMySQL()
		if err != nil {
			glog.ErrorF("connect %s error:%+v", dbname, err)
			return err
		}
		_Gorm.gormMaps.Store(dbname, orm)
		return nil
	}, 5, 3*time.Second)
	if err != nil {
		if oldGorm == nil {
			panic(err)
		}
		glog.WarnF("connect mysql failed, use the old one, database: %s error:%+v", dbname, err)
		return err
	}

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

func (i *ConnInfo) connectMySQL() (*gorm.DB, error) {
	if i.Dsn == "" {
		return nil, fmt.Errorf("in %s. DSN is empty", i.Dsn)
	}

	orm, err := gorm.Open(mysql.Open(i.Dsn), &gorm.Config{
		Logger:                 glog.NewDBLogger(i.Debug, time.Second),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	db, err := orm.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(i.MaxIdleConns)
	db.SetMaxOpenConns(i.MaxOpenConns)
	if i.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(i.ConnMaxIdleTime) * time.Second)
	}

	return orm, nil
}
