package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/v2/lib/middleware"
	"github.com/sunmi-OS/gocore/v2/lib/prometheus"
	tracing "github.com/sunmi-OS/gocore/v2/lib/tracing/gin/otel"
	"github.com/sunmi-OS/gocore/v2/utils/closes"
)

const (
	_HookShutdown hookType = "server_shutdown"
	_HookExit     hookType = "sys_exit"
)

type hookType string

type HookFunc func(c context.Context)

type GinEngine struct {
	Gin              *gin.Engine
	server           *http.Server
	timeout          time.Duration
	wg               sync.WaitGroup
	addrPort         string
	IgnoreReleaseLog bool
	hookMaps         map[hookType][]func(c context.Context)
}

func NewGinServer(ops ...Option) *GinEngine {
	cfg := defaultServerConfig
	for _, o := range ops {
		o(cfg)
	}

	g := gin.New()
	g.Use(logger(true), middleware.Recovery())
	if cfg.openTrace {
		// 引入链路追踪中间件
		endPointUrl := os.Getenv("ZIPKIN_BASE_URL")
		appName := os.Getenv("APP_NAME")
		if endPointUrl == "" || appName == "" {
			panic("开启链路追踪需要配置环境变量 ZIPKIN_BASE_URL 和 APP_NAME")
		}
		traceSampleRatio := os.Getenv("TRACE_SAMPLE_RATIO")
		sampleRatio := 1.0
		if traceSampleRatio != "" {
			sampleRatio = cast.ToFloat64(sampleRatio)
		}
		g.Use(tracing.ZipkinOtel(appName, endPointUrl, sampleRatio))
	}
	if !cfg.debug {
		gin.SetMode(gin.ReleaseMode)
	}
	// prometheus
	prometheus.NewPrometheus("app").Use(g)
	// default health check
	g.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	// pprof
	pp := g.Group("/debug/pprof")
	{
		pp.GET("/", gin.WrapF(pprof.Index))
		pp.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pp.GET("/profile", gin.WrapF(pprof.Profile))
		pp.GET("/symbol", gin.WrapF(pprof.Symbol))
		pp.GET("/trace", gin.WrapF(pprof.Trace))
		pp.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		pp.GET("/block", gin.WrapH(pprof.Handler("block")))
		pp.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		pp.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		pp.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		pp.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
	engine := &GinEngine{
		Gin:      g,
		addrPort: cfg.host + ":" + strconv.Itoa(cfg.port),
		server: &http.Server{
			Addr:         cfg.host + ":" + strconv.Itoa(cfg.port),
			Handler:      g.Handler(),
			ReadTimeout:  cfg.readTimeout,
			WriteTimeout: cfg.writeTimeout,
		},
		timeout:  cfg.readTimeout,
		wg:       sync.WaitGroup{},
		hookMaps: make(map[hookType][]func(c context.Context)),
	}
	return engine
}

func (g *GinEngine) Start() {
	// add common close hooks
	g.AddExitHook(func(c context.Context) {
		closes.Close()
	})

	// wait for signal
	go g.goNotifySignal()

	// start gin http server
	log.Printf("Listening and serving HTTP on %s\n", g.addrPort)
	if err := g.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("server.ListenAndServe(), error(%+v).", err))
		}
		log.Println("http: Server closed")
	}
	log.Println("waiting for process finished")
	// wait for process finished
	g.wg.Wait()
	log.Println("process exit")
}

// AddShutdownHook Add a hook function for when the GinServer service is shut down
func (g *GinEngine) AddShutdownHook(hooks ...HookFunc) *GinEngine {
	for _, fn := range hooks {
		if fn != nil {
			g.hookMaps[_HookShutdown] = append(g.hookMaps[_HookShutdown], fn)
		}
	}
	return g
}

// AddExitHook Add a hook function when the GinServer process exits
func (g *GinEngine) AddExitHook(hooks ...HookFunc) *GinEngine {
	for _, fn := range hooks {
		if fn != nil {
			g.hookMaps[_HookExit] = append(g.hookMaps[_HookExit], fn)
		}
	}
	return g
}

func (g *GinEngine) goNotifySignal() {
	g.wg.Add(1)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Printf("get a signal %s, stop the process\n", si.String())
			// close gin http server
			g.Close()
			// call before close hooks
			for _, fn := range g.hookMaps[_HookShutdown] {
				fn(context.Background())
			}
			// call after close hooks
			for _, fn := range g.hookMaps[_HookExit] {
				fn(context.Background())
			}
			// notify process exit
			g.wg.Done()
			runtime.Gosched()
			return
		case syscall.SIGHUP:
			log.Printf("get a signal %s\n", si.String())
		default:
			return
		}
	}
}

func (g *GinEngine) Close() {
	if g.server != nil {
		// disable keep-alives on existing connections
		g.server.SetKeepAlivesEnabled(false)
		_ = g.server.Shutdown(context.Background())
	}
}

// logger
func logger(ignoreRelease bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()
		if raw != "" {
			path = path + "?" + raw
		}

		// ignore logger output
		if gin.Mode() == gin.ReleaseMode && ignoreRelease {
			return
		}

		// End time
		end := time.Now()
		fmt.Fprintf(os.Stdout, "[GIN] %s | %3d | %13v | %15s | %-7s %#v\n%s", end.Format("2006/01/02 - 15:04:05"), c.Writer.Status(), end.Sub(start), c.ClientIP(), c.Request.Method, path, c.Errors.ByType(gin.ErrorTypePrivate).String())
	}
}
