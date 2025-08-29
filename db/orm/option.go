package orm

import (
	"os"
	"strconv"
	"strings"
)

type Options func(info *ConnInfo)

// ConnInfo connect info
type ConnInfo struct {
	Host            string // mysql's host
	Port            int    // mysql's port
	Username        string // the username to connect mysql
	Password        string // the password of username
	Database        string // the database for connect
	Dsn             string // data source name
	MultiStatements bool
	MaxIdleConns    int // Max of idle connects
	MaxOpenConns    int // max of open connects
	ConnMaxIdleTime int // max of connect idle time
	Debug           bool
	Type            string // the type of database
}

// NewConnInfo new connection info
// env format: {TYPE}_{DATABASE}_HOST, 比如 MYSQL_TEST_HOST
// database config item: DSN, HOST, PORT, USERNAME, PASSWORD, DATABASE
// if dsn not equals nil, use the dsn connect to db, or use host,port, username and password
func NewConnInfo(name string, opts ...Options) *ConnInfo {
	info := &ConnInfo{Type: "mysql"}
	for _, opt := range opts {
		opt(info)
	}
	if name == "" {
		name = info.Database
	}

	envPrefix := strings.ToUpper(info.Type + strings.TrimSpace(name) + "_")

	if info.Database == "" {
		info.Database = os.Getenv(envPrefix + "DATABASE")
	}

	if info.Dsn == "" {
		info.Dsn = os.Getenv(envPrefix + "DSN")
	}
	if info.Host == "" {
		info.Host = os.Getenv(envPrefix + "HOST")
	}
	if info.Port == 0 {
		info.Port, _ = strconv.Atoi(os.Getenv(envPrefix + "PORT"))
	}
	if info.Username == "" {
		info.Username = os.Getenv(envPrefix + "USERNAME")
	}
	if info.Password == "" {
		info.Password = os.Getenv(envPrefix + "PASSWORD")
	}

	return info
}

func WithDebug(debug bool) Options {
	return func(info *ConnInfo) {
		info.Debug = debug
	}
}

func WithDsn(dsn string) func(options *ConnInfo) {
	return func(options *ConnInfo) {
		options.Dsn = dsn
	}
}

func WithHost(host string) Options {
	return func(c *ConnInfo) {
		c.Host = host
	}
}

func WithPort(port int) Options {
	return func(c *ConnInfo) {
		c.Port = port
	}
}

func WithUsername(username string) Options {
	return func(c *ConnInfo) {
		c.Username = username
	}
}

func WithPassword(password string) Options {
	return func(c *ConnInfo) {
		c.Password = password
	}
}
func WithDatabase(database string) Options {
	return func(c *ConnInfo) {
		c.Database = database
	}
}

func WithMaxIdleConns(maxIdleConns int) Options {
	return func(c *ConnInfo) {
		c.MaxIdleConns = maxIdleConns
	}
}

func WithMaxOpenConns(maxOpenConns int) Options {
	return func(c *ConnInfo) {
		c.MaxOpenConns = maxOpenConns
	}
}

func WithConnMaxIdleTime(connMaxIdleTime int) Options {
	return func(c *ConnInfo) {
		c.ConnMaxIdleTime = connMaxIdleTime
	}
}
