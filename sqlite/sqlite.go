package sqlite

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Gorm sync.Map

// 初始化Gorm
func NewDB(dbname string) {

	var orm *gorm.DB
	var err error

	for orm, err = gorm.Open("sqlite3", "./"+dbname+".db"); err != nil; {
		fmt.Println("Database connection exception! 5 seconds to retry")
		time.Sleep(5 * time.Second)
		orm, err = gorm.Open("sqlite3", "./"+dbname+".db")
	}
	Gorm.LoadOrStore(dbname, orm)
}

// 通过名称获取Gorm实例
func GetORMByName(dbname string) *gorm.DB {

	v, _ := Gorm.Load(dbname)
	return v.(*gorm.DB)
}

// 获取默认的Gorm实例
func GetORM() *gorm.DB {

	v, _ := Gorm.Load("dbDefault")
	return v.(*gorm.DB)
}
