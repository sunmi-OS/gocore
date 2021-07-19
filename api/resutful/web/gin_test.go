package web

import (
	"testing"

	"github.com/sunmi-OS/gocore/v2/api/ecode"

	"github.com/gin-gonic/gin"
)

var banner = `
                                   _
 ___   _   _   _ __    _ __ ___   (_)
/ __| | | | | | '_ \  | '_ * _ \  | |
\__ \ | |_| | | | | | | | | | | | | |
|___/  \__,_| |_| |_| |_| |_| |_| |_|
`

func TestInitGin(t *testing.T) {
	//c := &Config{Port: ":2233"}
	//g := InitGin(c).Release()
	//
	//initRouteG(g.Gin)
	//
	//g.Start(banner)
	//
	//ch := make(chan os.Signal)
	//signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	//for {
	//	si := <-ch
	//	switch si {
	//	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
	//		time.Sleep(time.Second)
	//		// todo something
	//
	//		time.Sleep(time.Second)
	//		return
	//	case syscall.SIGHUP:
	//	default:
	//		return
	//	}
	//}
}

func initRouteG(g *gin.Engine) {
	g.GET("/gin/ping", func(c *gin.Context) {
		e := ecode.New(2233, "success")
		JSON(c, nil, e)
	})

	g.GET("/gin/file", func(c *gin.Context) {
		//err := ecode.New(2323, "asdsda")
		File(c, "echo_test.go", "gin.go")
	})

	g.GET("/gin/page", func(c *gin.Context) {
		e := ecode.New(2233, "success")

		rsp := []string{"1", "2", "3", "4", "5"}
		JSON(c, Pager{PageNo: 1, PageSize: 5}.Apply(10, rsp), e)
	})
}
