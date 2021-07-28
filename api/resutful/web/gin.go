package web

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	Gin       *gin.Engine
	isRelease bool
	addr      string
}

func InitGin(c *Config) *GinEngine {
	g := gin.Default()
	if !strings.Contains(strings.TrimSpace(c.Port), ":") {
		c.Port = ":" + c.Port
	}
	engine := &GinEngine{Gin: g, addr: c.Host + c.Port}
	return engine
}

// Release release
func (g *GinEngine) Release() *GinEngine {
	gin.SetMode(gin.ReleaseMode)
	return g
}

func (g *GinEngine) Start() {
	go func() {
		if err := g.Gin.Run(g.addr); err != nil {
			panic(fmt.Sprintf("web server port(%s) run error(%+v).", g.addr, err))
		}
	}()
}
