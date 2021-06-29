package gormv2

import (
	"database/sql"

	"github.com/iGoogle-ink/gopher/xtime"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrNoRow          = sql.ErrNoRows
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

// MySQLConfig mysql config.
type MySQLConfig struct {
	DSN            string          // data source name.
	MaxOpenConn    int             // pool, e.g:10
	MaxIdleConn    int             // pool, e.g:100
	MaxConnTimeout xtime.Duration  // connect max life time. Unmarshal config file e.g: 10s、2m、1m10s
	LogLevel       logger.LogLevel // Silent、Info、Warn、Error, default Warn
	SlowThreshold  xtime.Duration  // slow sql log. Unmarshal config file e.g: 100ms、200ms、300ms、1s, default 200ms
}
