package xxljob

type Option struct {
	AppName     string   // xxl-job executor's app name
	Port        string   // the port of xxl-job executor, default is 9999
	LogLevel    LogLevel // log level, example: debug、info、warn、error, default is info
	accessToken string   // access token of xxl-job executor
	serverAddr  string   // the address of xxl-job admin, if empty, will use env XXL_JOB_SERVER_ADDR
}

func (c *Option) WithAccessToken(token string) *Option {
	c.accessToken = token
	return c
}

func (c *Option) WithServerAddr(addr string) *Option {
	c.serverAddr = addr
	return c
}
