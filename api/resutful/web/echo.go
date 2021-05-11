package web

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

type EchoEngine struct {
	Echo *echo.Echo
	addr string
}

func InitEcho(c *Config) *EchoEngine {
	e := echo.New()

	if !strings.Contains(strings.TrimSpace(c.Port), ":") {
		c.Port = ":" + c.Port
	}

	engine := &EchoEngine{Echo: e, addr: c.Host + c.Port}

	e.Use(engine.Logger())
	e.Use(engine.Recover())
	return engine
}

// Release
func (e *EchoEngine) Release() *EchoEngine {
	e.Echo.Debug = false
	e.Echo.HideBanner = true
	e.Echo.HidePort = true
	return e
}

func (e *EchoEngine) Start(banner ...string) {
	if len(banner) > 0 {
		fmt.Println(banner[0])
	}
	go func() {
		if err := e.Echo.Start(e.addr); err != nil {
			panic(fmt.Sprintf("web server port(%s) run error(%+v).", e.addr, err))
		}
	}()
}
