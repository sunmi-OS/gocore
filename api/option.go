package api

import "time"

type Config struct {
	host         string        // host
	port         int           // port default 2233
	readTimeout  time.Duration // read_timeout, default 60 (s)
	writeTimeout time.Duration // write_timeout, default 60 (s)
	debug        bool          // is show log
	openTrace    bool          // is open trace
}

type Option func(o *Config)

var defaultServerConfig = &Config{
	host:         "",
	port:         2233,
	readTimeout:  60,
	writeTimeout: 60,
	debug:        false,
	openTrace:    false,
}

// ServerHost 设置host
func ServerHost(addr string) Option {
	return func(c *Config) { c.host = addr }
}

// ServerPort 设置端口
func ServerPort(port int) Option {
	return func(c *Config) { c.port = port }
}

// ServerTimeout 设置超时时间
func ServerTimeout(dur time.Duration) Option {
	return func(o *Config) {
		o.readTimeout = dur
		o.writeTimeout = dur
	}
}

// ServerDebug 设置超时时间
func ServerDebug(debug bool) Option {
	return func(o *Config) { o.debug = debug }
}

// OpenTrace 设置超时时间
func OpenTrace(open bool) Option {
	return func(o *Config) { o.openTrace = open }
}
