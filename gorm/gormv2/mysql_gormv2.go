package gormv2

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGormV2(c *MySQLConfig) (db *gorm.DB) {
	lc := logger.Config{
		SlowThreshold: 200 * time.Millisecond, // 慢 SQL 阈值
		LogLevel:      logger.Warn,            // Log level
		Colorful:      false,                  // 禁用彩色打印，日志平台会打印出颜色码，影响日志观察
	}
	if c.LogLevel != 0 {
		lc.LogLevel = c.LogLevel
	}
	if c.SlowThreshold != 0 {
		lc.SlowThreshold = time.Duration(c.SlowThreshold)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "[GORM] >> ", 64|log.Ldate|log.Lmicroseconds), // io writer
		lc,
	)
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database error:%+v", err))
	}
	sql, err := db.DB()
	if err != nil {
		panic(err)
	}
	sql.SetMaxIdleConns(c.MaxIdleConn)
	sql.SetMaxOpenConns(c.MaxOpenConn)
	sql.SetConnMaxLifetime(time.Duration(c.MaxConnTimeout))
	return db
}
