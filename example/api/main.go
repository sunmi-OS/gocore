package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"sort"
	"os"
	"github.com/urfave/cli"
	"github.com/sunmi-OS/gocore/api"
	"github.com/sunmi-OS/gocore/viper"
)

type EchoApi struct {
}

var eApi EchoApi

func (a *EchoApi) echoStart(c *cli.Context) error {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.POST("/", func(c echo.Context) error {

		request := api.NewRequest(c)
		response := api.NewResponse(c)

		err := request.InitDES()
		if err != nil {
			return response.RetError(err, -1)
		}

		msg := request.DESParam("msg").GetString()

		return response.RetSuccess(msg)
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
	return nil
}

func main() {

	app := cli.NewApp()
	app.Name = "IOT-seanbox"
	app.Usage = "IOT-seanbox"
	app.Email = "wenzhenxi@sunmi.com"
	app.Version = "1.0.0"
	app.Usage = "IOT-seanbox"

	// 初始化配置
	viper.NewConfig("config", "conf")

	// 指定对于的命令
	app.Commands = []cli.Command{
		{
			Name:    "api",
			Aliases: []string{"a"},
			Usage:   "api",
			Subcommands: []cli.Command{
				{
					Name:   "start",
					Usage:  "开启API-DEMO",
					Action: eApi.echoStart,
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)

}
