package web

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	Gin  *gin.Engine
	port string
}

func InitGin(c *Config) *GinEngine {
	g := gin.Default()

	engine := &GinEngine{Gin: g, port: c.Port}
	if !strings.Contains(strings.TrimSpace(c.Port), ":") {
		engine.port = ":" + c.Port
	}

	g.Use(engine.CORS())
	g.Use(engine.Recovery())
	return engine
}

func (g *GinEngine) Release() *GinEngine {
	gin.SetMode(gin.ReleaseMode)
	return g
}

func (g *GinEngine) Start(banner ...string) {
	if len(banner) > 0 {
		fmt.Println(banner[0])
	}
	go func() {
		if err := g.Gin.Run(g.port); err != nil {
			panic(fmt.Sprintf("web server port(%s) run error(%+v).", g.port, err))
		}
	}()
}
