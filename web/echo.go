package web

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

type EchoEngine struct {
	Echo *echo.Echo
	port string
}

func InitEcho(c *Config) *EchoEngine {
	e := echo.New()

	engine := &EchoEngine{Echo: e, port: c.Port}
	if !strings.Contains(strings.TrimSpace(c.Port), ":") {
		engine.port = ":" + c.Port
	}

	e.Use(engine.Logger())
	e.Use(engine.Recover())
	return engine
}

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
		if err := e.Echo.Start(e.port); err != nil {
			panic(fmt.Sprintf("web server port(%s) run error(%+v).", e.port, err))
		}
	}()
}
