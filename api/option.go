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
	port:         8080,
	readTimeout:  60 * time.Second,
	writeTimeout: 60 * time.Second,
	debug:        false,
	openTrace:    true,
}

// WithServerHost 设置host
func WithServerHost(addr string) Option {
	return func(c *Config) { c.host = addr }
}

// WithServerPort 设置端口
func WithServerPort(port int) Option {
	return func(c *Config) { c.port = port }
}

// WithServerTimeout 设置超时时间
func WithServerTimeout(dur time.Duration) Option {
	return func(o *Config) {
		o.readTimeout = dur
		o.writeTimeout = dur
	}
}

// WithServerDebug 设置超时时间
func WithServerDebug(debug bool) Option {
	return func(o *Config) { o.debug = debug }
}

// WithOpenTrace 设置超时时间
func WithOpenTrace(open bool) Option {
	return func(o *Config) { o.openTrace = open }
}
