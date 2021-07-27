package web

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/utils/limit"
	"github.com/sunmi-OS/gocore/v2/utils/trace"
)

type GinEngine struct {
	Gin       *gin.Engine
	Tracer    *trace.Tracer
	isRelease bool
	addr      string
}

func InitGin(c *Config) *GinEngine {
	g := gin.Default()

	if c.Limit != nil && c.Limit.Rate != 0 {
		g.Use(limit.NewLimiter(c.Limit).GinLimit())
	}

	if !strings.Contains(strings.TrimSpace(c.Port), ":") {
		c.Port = ":" + c.Port
	}

	engine := &GinEngine{Gin: g, addr: c.Host + c.Port}
	if c.Trace != nil {
		engine.Tracer = trace.NewTracer(c.Trace)
		g.Use(engine.Tracer.GinTrace())
	}
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

func (g *GinEngine) Close() {
	if g.Tracer != nil {
		g.Tracer.Close()
	}
}
