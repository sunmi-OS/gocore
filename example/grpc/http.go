package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/sunmi-OS/gocore/log"
	"golang-example/pkg/istio"
	printpb "golang-example/pkg/proto/print"
	"net/http"
	"time"
)

func main() {

	e := echo.New()

	printpb.Init(":8080", 3000)

	e.Any("/grpc-test", func(c echo.Context) error {

		trace := istio.SetHttp(c.Request().Header)
		req := &printpb.Request{
			Message: "test",
		}
		resp, err := printpb.PrintOk(req, trace)
		if err != nil {
			log.Sugar.Errorf("请求失败: %s", err.Error())
			return err
		}
		fmt.Println(resp.Data)

		return c.JSON(200, resp.Data)
	})

	s := &http.Server{
		Addr:         ":80",
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}

	e.StartServer(s)
}
